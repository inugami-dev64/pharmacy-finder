import { pushState } from "$app/navigation";
import type { HttpError } from "$lib/http-error";

export const PAGER_LIMIT = 10

export class PharmacyReview {
    id: number | undefined;
    prescriptionType: string | undefined;
    stars: number | undefined;
    hrtKind: string | undefined;
    nationality: string | null | undefined;
    review: string | null | undefined;
    modCode: string | undefined;
    createdAt: number | undefined;
    updatedAt: number | undefined;

    // TODO: Fix error handling because right now exceptions are caught in data methods

    /**
     * Retrieve an array of PharmacyReviews from backend
     *
     * @param id ID of the pharmacy whose reviews to query
     * @param key pager key value. In this case the updatedAt timestamp in unix millis.
     * @param uniqueKey pager unique key values. In this case, the ID of the last review
     * @returns a promise to PharmacyReview array (meow :3)
     */
    public static async readReviews(id: number, key: number | undefined, uniqueKey: number | undefined): Promise<PharmacyReview[]> {
        return await fetch(`/api/v1/pharmacies/${id}/reviews?l=${PAGER_LIMIT}${key != null ? `&k=${key}` : ""}${uniqueKey != null ? `&uk=${uniqueKey}` : ""}`)
            .then(async res => {
                if (res.status != 200) {
                    let err: HttpError = await res.json();
                    console.log(err);
                    throw new Error(`Failed to fetch pharmacy reviews for pharmacy ${id}`)
                }

                let data: PharmacyReview[] = await res.json();
                return data
            })
            .catch(e => {
                console.log(e);
                return [];
            })
    }

    /**
     * Create a new pharmacy review from current instance and submit it to backend
     *
     * @param id ID of the pharmacy whose review is going to be submitted
     * @returns a promise to created PharmacyReview instance
     */
    public async createReview(id: number): Promise<PharmacyReview> {
        // undefine fields which we do not want to submit
        let copy = new PharmacyReview;
        Object.assign(copy, this);

        copy.id = undefined;
        copy.modCode = undefined;
        copy.createdAt = undefined;
        copy.updatedAt = undefined;

        return await fetch(`/api/v1/pharmacies/${id}/reviews`, {
            method: "POST",
            body: JSON.stringify(copy),
            headers: {
                "Content-Type": "application/json"
            }
        })
            .then(async res => {
                if (res.status != 201) {
                    let err: HttpError = await res.json();
                    console.log(err);
                    throw new Error(`Failed to create a new post`);
                }

                let data: PharmacyReview = await res.json();
                return data;
            })
            .catch(e => {
                console.log(e);
                return new PharmacyReview;
            });
    }

    /**
     * Update the pharmacy review
     *
     * @param id specifies the pharmacy ID whose review to update
     * @returns an updated PharmacyReview instance
     */
    public async updateReview(id: number): Promise<PharmacyReview> {
        // undefine fields which we do not want to submit
        let copy = new PharmacyReview;
        Object.assign(copy, this);

        copy.id = undefined;
        copy.createdAt = undefined;
        copy.updatedAt = undefined;

        return await fetch(`/api/v1/pharmacies/${id}/reviews/${this.id}`, {
            method: "PATCH",
            body: JSON.stringify(copy),
            headers: {
                "Content-Type": "application/json"
            }
        })
            .then(async res => {
                if (res.status != 200) {
                    let err: HttpError = await res.json();
                    console.log(err);
                    throw new Error(`Failed to create a new post`);
                }

                let data: PharmacyReview = await res.json();
                return data;
            })
            .catch(e => {
                console.log(e);
                return new PharmacyReview;
            })
    }
}