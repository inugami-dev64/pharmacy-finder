import type { HttpError } from "$lib/http-error";

export class PharmacyInfo {
    id: number | undefined;
    chain: string | undefined;
    name: string | undefined;
    address: string | undefined;
    city: string | undefined
    county: string | undefined
    postalCode: number | undefined;
    phoneNumber: string | undefined
    lat: number | undefined;
    lng: number | undefined;
}

export async function getPharmacies(): Promise<PharmacyInfo[]> {
    return await fetch("/api/v1/pharmacies")
        .then(async res => {
            if (res.status != 200) {
                let err: HttpError = await res.json();
                console.log(err);
                throw new Error(`Failed to fetch pharmacies: ${err.msg}`);
            }

            let data: PharmacyInfo[] = await res.json()
            return data
        })
        .catch(e => {
            console.log(e);
            return [];
        }) /* TODO: Implement a better error signaling system so that the user can see errors */
}