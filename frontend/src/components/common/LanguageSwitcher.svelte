<script lang="ts">
    import { locale } from 'svelte-i18n';

    let languageSelections = $state([
        {
            code: "en",
            language: "🇬🇧 English"
        },
        {
            code: "et",
            language: "🇪🇪 Eesti"
        }
    ])

    console.log("en")
    let selectedValue: string | undefined = $state()

    function changeLanguage(e: Event) {
        e.preventDefault()
        locale.set(selectedValue);
    }
</script>

<label for="langSelect"></label>
<select
    bind:value={selectedValue}
    onchange={changeLanguage}
    name="langSelect"
>
    {#each languageSelections as selection}
    <option value="{selection.code}">{selection.language}</option>
    {/each}
</select>

<style>
    select {
        appearance: none;
        padding: 0.3em 0.75em;
        border: 1px solid var(--primary-text-color);
        font-family: inherit;
        font-size: inherit;
        line-height: inherit;
        background-color: rgba(253, 206, 206, 0.3);
        width: fit-content;
        border-radius: 0.5em;
        margin: 0.5em;
    }

    select:hover {
        cursor: pointer;
    }
</style>