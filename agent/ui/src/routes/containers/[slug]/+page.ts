export const load = ({ params }) => {
  console.log("hit load");
  return {
    slug: params.slug,
  };
};
