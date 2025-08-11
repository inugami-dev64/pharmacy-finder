
export class PharmacyInfo {
    chain: string | undefined;
    name: string | undefined;
    aadress: string | undefined;
    postalCode: number | undefined;
    city: string | undefined
    county: string | undefined
    phoneNumber: string | undefined
    latitude: number | undefined;
    longitude: number | undefined;
    avgRating: number | undefined;
}

export async function getPharmacies(): Promise<PharmacyInfo[]> {
    return [
        {
            chain: "Benu",
            name: "AIA APTEEK",
            aadress: "Narva mnt 7",
            postalCode: 10017,
            city: "Tallinn",
            county: "Harjumaa",
            phoneNumber: "+3726109490",
            latitude: 59.437264,
            longitude: 24.760051,
            avgRating: 1
        }
    ]
}