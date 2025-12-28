/**
 * Here we define all possible languages to use
 */

export interface Language {
    code: string
    language: string
}

export const languages: Language[] = [
    {
        code: "en",
        language: "English"
    },
    {
        code: "et",
        language: "eesti"
    }
]