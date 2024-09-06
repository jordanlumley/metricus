package agent

import (
	"fmt"
	"net/http"

	metricus "github.com/jordanlumley/metricus/sdk"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func init() {
	// output := zerolog.ConsoleWriter{
	// 	Out:        os.Stderr,
	// 	TimeFormat: zerolog.TimeFormatUnix,
	// }

	// zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// log.Logger = log.Output(output)
}

func StartAPI() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("method", v.Method).
				Err(v.Error).
				Msg("REQUEST")

			return nil
		},
	}))

	dockerService, err := metricus.NewDockerService()
	if err != nil {
		log.Fatal().Err(err).Msg("error creating docker service")
	}

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	v1 := e.Group("/api/v1")
	v1.GET("/metrics/events", func(c echo.Context) error {
		UpgradeSSE(c.Response())

		metricsStream := make(chan []byte)
		ctx := c.Request().Context()

		go func() {
			defer close(metricsStream)

			// start scraping the metrics endpoint every 10 seconds
			agent := New(Options{
				ScrapeTargets: []ScrapeOptions{
					{
						Host:            "http://example_client:8080",
						IntervalSeconds: 2,
						Channel:         metricsStream,
					},
				},
			})
			agent.Start()
		}()

		for {
			select {
			case <-ctx.Done():
				return nil
			case message, ok := <-metricsStream:
				if !ok {
					return nil
				}

				if err := SendSSE(c.Response(), message); err != nil {
					log.Error().Err(err).Msg("error sending metricsStream message")
					return err
				}
			}
		}
	})

	v1.GET("/containers", func(c echo.Context) error {
		containers, err := dockerService.GetContainers(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error getting containers")
		}

		return c.JSON(http.StatusOK, containers)
	})

	v1.GET("/containers/:containerId", func(c echo.Context) error {
		containerID := c.Param("containerId")
		container, err := dockerService.GetContainer(c.Request().Context(), containerID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error getting container(%s)", containerID))
		}

		return c.JSON(http.StatusOK, container)
	})

	v1.GET("/containers/:containerId/stats", func(c echo.Context) error {
		containerID := c.Param("containerId")
		stats, err := dockerService.GetContainerMetrics(c.Request().Context(), containerID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error getting container(%s)", containerID))
		}

		return c.JSON(http.StatusOK, stats)
	})

	v1.GET("/containers/:containerId/stats/test", func(c echo.Context) error {
		containerID := c.Param("containerId")

		UpgradeSSE(c.Response())

		metricsStream := make(chan []byte)
		ctx := c.Request().Context()

		go func() {
			defer close(metricsStream)
			if err := dockerService.StreamContainerMetrics(c.Request().Context(), containerID, metricsStream); err != nil {
				log.Error().Err(err).Msg("error streaming metrics")
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return nil
			case message, ok := <-metricsStream:
				fmt.Println("sending message metrics", containerID)
				if !ok {
					return nil
				}

				if err := SendSSE(c.Response(), message); err != nil {
					log.Error().Err(err).Msg("error sending metricsStream message")
					return err
				}
			}
		}
	})

	v1.GET("/containers/:containerId/logs", func(c echo.Context) error {
		containerID := c.Param("containerId")
		logs, err := dockerService.GetLogs(c.Request().Context(), containerID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error getting container(%s) logs", containerID))
		}

		return c.JSON(http.StatusOK, logs)
	})

	v1.GET("/containers/:containerId/logs/events", func(c echo.Context) error {
		containerID := c.Param("containerId")

		UpgradeSSE(c.Response())

		logStream := make(chan []byte)
		ctx := c.Request().Context()

		go func() {
			defer close(logStream)
			if err := dockerService.StreamLogs(c.Request().Context(), containerID, logStream); err != nil {
				log.Error().Err(err).Msg("error streaming logs")
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return nil
			case message, ok := <-logStream:
				fmt.Println("sending message logs", containerID)
				if !ok {
					return nil
				}

				if err := SendSSE(c.Response(), message); err != nil {
					log.Error().Err(err).Msg("error sending logStream message")
					return err
				}
			}
		}
	})

	// Start the server
	if err = e.Start(":8888"); err != nil {
		log.Fatal().Err(err).Msg("error starting api server")
	}
}
