<script lang="ts">
    import { localeSwitcher } from "$lib/service/locale";
    import LanguageButton from "../buttons/LanguageButton.svelte";
    import DropdownMenuContainer from "./DropdownMenuContainer.svelte";

    export let size: number;
    export let title: string | undefined = "Change language";
    export let color: string | undefined = "#999";
</script>

<DropdownMenuContainer
    size={size}
    ButtonComponent={LanguageButton}
    title={title}
    color={color}
>
    <button onclick={() => {
        localeSwitcher.setLocale("en");
    }}>
        <span style="font-size: 32px;">🇬🇧</span>
        <p>English</p>
    </button>
    <button onclick={() => {
        localeSwitcher.setLocale("et");
    }}>
        <span style="font-size: 32px;">🇪🇪</span>
        <p>eesti</p>
    </button>
</DropdownMenuContainer>