import type { HttpError } from "$lib/http-error";
import { Point } from "$lib/utils/point";

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

/**
 * Tier ratings are mainly used in the tier-list
 */
export class PharmacyTierRating {
    id: number | undefined;
    name: string | undefined;
    avgRating: number | undefined;
    avgERating: number | undefined;
    avgTRating: number | undefined;

    public static async readPharmacyTierRatings(sw?: Point, ne?: Point): Promise<PharmacyTierRating[]> {
        if (sw == null)
            sw = new Point(-90, -90);
        if (ne == null)
            ne = new Point(90, 90);

        return await fetch(`/api/v1/pharmacies/ratings?sw=${sw.lat},${sw.lng}&ne=${ne.lat},${ne.lng}`)
            .then(async res => {
                if (res.status != 200) {
                    let err: HttpError = await res.json();
                    console.log(err);
                    throw new Error(`Failed to fetch pharmacy tier ratings`);
                }

                let data: PharmacyTierRating[] = await res.json();
                return data;
            })
            .catch(e => {
                console.log(e);
                return [];
            })
    }
}