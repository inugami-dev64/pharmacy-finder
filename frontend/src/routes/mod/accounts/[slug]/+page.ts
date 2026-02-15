export function load({ params }) {
    const userId: string = params.slug;
    
    return {
        slug: params.slug
    }
}