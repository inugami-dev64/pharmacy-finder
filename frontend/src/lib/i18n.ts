import { _, getLocaleFromNavigator } from "svelte-i18n";

import { addMessages, init } from "svelte-i18n";

import en from "../lang/en.json";
import et from "../lang/et.json";

addMessages("en", en);
addMessages("et", et);

const initialLocale = getLocaleFromNavigator();

init({
    initialLocale: initialLocale,
    fallbackLocale: "en"
});