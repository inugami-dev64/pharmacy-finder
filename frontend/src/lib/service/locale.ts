import { locale } from "svelte-i18n";

const LOCAL_STORAGE_KEY = "LOCALE";

export class LocaleSwitcher {
    /**
     * Retrieves and returns the current locale value to use
     *
     * @returns The current locale (default fallback locale: en)
     */
    public getLocale(): string {
        const localeValue = localStorage.getItem(LOCAL_STORAGE_KEY);
        if (localeValue == null)
            return "en";

        return localeValue;
    }

    /**
     * Sets default locale for the page
     */
    public setDefault() {
        const localeValue = localStorage.getItem(LOCAL_STORAGE_KEY);
        if (localeValue == null)
            locale.set("en");
        else locale.set(localeValue);
    }

    /**
     * Modifies the current locale of the webpage
     *
     * @param locale specifies the locale to change to
     */
    public setLocale(localeValue: string) {
        locale.set(localeValue);
        localStorage.setItem(LOCAL_STORAGE_KEY, localeValue);
    }
}

export const localeSwitcher = new LocaleSwitcher();