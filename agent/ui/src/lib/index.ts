// place files you want to import through the `$lib` alias in this folder.

export async function fetchContainers() {
  try {
    const response = await fetch("http://localhost:8888/api/v1/containers");
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Error fetching containers:", error);
    throw error;
  }
}

export async function fetchContainer(container: String) {
  try {
    const response = await fetch(
      `http://localhost:8888/api/v1/containers/${container}`
    );
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Error fetching container:", error);
    throw error;
  }
}

export async function fetchContainerStats(container: String) {
  try {
    const response = await fetch(
      `http://localhost:8888/api/v1/containers/${container}/stats`
    );
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Error fetching container stats:", error);
    throw error;
  }
}

export async function fetchMetrics() {
  try {
    const response = await fetch(`http://localhost:8888/api/v1/metrics`);
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Error fetching metrics:", error);
    throw error;
  }
}
