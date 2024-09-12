export const load = ({ params }: any) => {
  console.log("hit load");
  return {
    slug: params.slug,
  };
};
