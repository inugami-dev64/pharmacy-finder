
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
    overallRating: number | undefined;
    acceptanceRating: number | undefined;
    eRating: number | undefined;
    tRating: number | undefined;
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
            overallRating: 3.8,
            acceptanceRating: 2.5,
            eRating: 3.0,
            tRating: 1
        }
    ]
}