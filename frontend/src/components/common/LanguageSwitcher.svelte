<script lang="ts">
    import { localeSwitcher } from '$lib/service/locale';
    import { languages } from '$lib/utils/languages';
    import { onMount } from 'svelte';

    let selectedValue: string | undefined = $state();

    onMount(() => selectedValue = localeSwitcher.getLocale());

    function changeLanguage(e: Event) {
        e.preventDefault();
        localeSwitcher.setLocale(selectedValue || "");
    }
</script>

<label for="langSelect"></label>
<select
    bind:value={selectedValue}
    onchange={changeLanguage}
    name="langSelect"
>
    {#each languages as selection}
    <option value="{selection.code}" selected={selection.language === selectedValue}>{selection.language}</option>
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