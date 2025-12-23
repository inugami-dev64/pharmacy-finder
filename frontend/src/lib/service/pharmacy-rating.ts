import type { HttpError } from "$lib/http-error";

export class PharmacyRating {
    id: number | undefined;
    stars: number | undefined;
    hrtKind: string | undefined;

    /**
     * Retrieve pharmacy ratings (aggregated values) for provided pharmacy
     *
     * @param id ID of the pharmacy whose ratings to query for
     * @returns a promise to PharmacyRating array
     */
    public static async readPharmacyRatings(id: number): Promise<PharmacyRating[]> {
        return await fetch(`/api/v1/pharmacies/${id}/ratings`)
            .then(async res => {
                if (res.status != 200) {
                    let err: HttpError = await res.json();
                    console.log(err);
                    throw new Error(`Failed to fetch pharmacy ratings for pharmacy ${id}`);
                }

                let data: PharmacyRating[] = await res.json()
                return data
            })
            .catch(e => {
                console.log(e);
                return [];
            })
    }
}