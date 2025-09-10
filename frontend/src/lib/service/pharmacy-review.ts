import type { HttpError } from "$lib/http-error";

export class PharmacyReview {
    id: number | undefined;
    prescriptionType: string | undefined;
    stars: number | undefined;
    hrtKind: string | undefined;
    nationality: string | null | undefined;
    review: string | null | undefined;
    createdAt: number | undefined;
    updatedAt: number | undefined;
}

export async function getPharmacyReviews(id: number, key: number|null, uniqueKey: number|null): Promise<PharmacyReview[]> {
    return await fetch(`/api/v1/pharmacies/${id}/reviews?l=20${key != null ? `&k=${key}` : ""}${uniqueKey != null ? `&uk=${uniqueKey}` : ""}`)
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