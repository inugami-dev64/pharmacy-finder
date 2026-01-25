<script lang="ts">
    import { PharmacyInfo } from "$lib/service/pharmacy-info";
    import { searchViewZIndex } from "$lib/utils/z-indices";
    import CloseButton from "../../common/icons/buttons/CloseButton.svelte";
    import SearchIcon from "../../common/icons/SearchIcon.svelte";
    import Overlay from "../../common/Overlay.svelte";


    // Pharmacy icons
    import ApothekaMarker from "$lib/assets/markers/apotheka.png";
    import SudameapteekMarker from "$lib/assets/markers/sudameapteek.png";
    import BenuMarker from "$lib/assets/markers/benu.png";
    import EuroapteekMarker from "$lib/assets/markers/euroapteek.png";

    let {
        onSelect,
        onClose,
        pharmacies,
        minWidth = 360,
        minHeight = 500
    }: {
        onSelect: (pharmacy: PharmacyInfo) => void,
        onClose: () => void,
        pharmacies: PharmacyInfo[],
        minWidth?: number,
        minHeight?: number
    } = $props();

    // stateful values
    let filteredPharmacies: PharmacyInfo[] = $state([]);
    let currentKeyBoardFocusIndex: number = -1;

    function onInput(e: Event & { currentTarget: EventTarget & HTMLInputElement; }) {
        e.preventDefault()
        const tokens = (e.target as HTMLInputElement).value.split(" ");

        if (tokens.length === 0 || tokens.length === 1 && tokens[0] === "") {
            filteredPharmacies = [];
            return;
        }

        filteredPharmacies = pharmacies.filter(v => {
            const searchStr = v.name?.toLowerCase() ?? "";
            let containsAll = true;
            for (const token of tokens) {
                containsAll = containsAll && searchStr.includes(token.toLowerCase());
            }

            return containsAll;
        }).sort((v1, v2) => {
            if ((v1.name ?? "") < (v2.name ?? ""))
                return -1;
            else if (v1.name === v2.name)
                return 0;
            return 1;
        }).slice(0, 5);

    }

    // When SearchModal is rendered, make the searchText visible
    let searchText: HTMLInputElement;
    $effect(() => searchText.focus());

    // Suggestions container
    let suggestions: HTMLDivElement;

    function onKeyDownEvent(e: KeyboardEvent & {currentTarget: EventTarget & Window}) {
        if (e.key == "Escape") {
            onClose();
        } else if (e.key == "ArrowDown") {
            currentKeyBoardFocusIndex++;
        } else if (e.key == "ArrowUp") {
            currentKeyBoardFocusIndex--;
        }

        // clip values
        if (currentKeyBoardFocusIndex < -1)
            currentKeyBoardFocusIndex = -1
        if (currentKeyBoardFocusIndex >= suggestions.children.length)
            currentKeyBoardFocusIndex = suggestions.children.length - 1;

        if (currentKeyBoardFocusIndex == -1)
            searchText.focus();
        else
            (suggestions.children.item(currentKeyBoardFocusIndex) as HTMLButtonElement).focus();
    }
</script>

<Overlay zIndex={searchViewZIndex}>
    <div class="search-modal" style="--minWidth: {minWidth}px; --minHeight: {minHeight}">
        <div class="search-box">
            <SearchIcon size={32}/>
            <input type="text" placeholder="Search" oninput={onInput} bind:this={searchText}/>
            <CloseButton size={32} on:click={onClose}/>
        </div>

        <div bind:this={suggestions}>
            {#each filteredPharmacies as pharmacy}
                <button type="button" class="suggestion" onclick={(_) => onSelect(pharmacy)} tabindex="0">
                    {#if pharmacy.chain === "Apotheka"}
                        <img src="{ApothekaMarker}" alt="Apotheka">
                    {:else if pharmacy.chain === "Südameapteek"}
                        <img src="{SudameapteekMarker}" alt="Südameapteek">
                    {:else if pharmacy.chain === "Benu"}
                        <img src="{BenuMarker}" alt="Benu">
                    {:else if pharmacy.chain === "Euroapteek"}
                        <img src="{EuroapteekMarker}" alt="Euroapteek">
                    {/if}
                    <div class="pharmacy-name">
                        <h3>{pharmacy.name}</h3>
                    </div>
                </button>
            {/each}
        </div>
    </div>
</Overlay>

<svelte:window on:keydown|stopPropagation={onKeyDownEvent}/>

<style>
    .search-modal {
        display: flex;
        flex-direction: column;
        background-color: rgba(0, 0, 0, 0);
        width: 50%;
        height: 50%;
        min-width: var(--minWidth);
        min-height: var(--minHeight);
    }

    .search-box {
        display: flex;
        align-items: center;
        padding-left: 0.5em;
        background-color: #ffffff;
        box-sizing: border-box;
        flex-direction: row;
        align-items: center;
        height: 48px;
        justify-content: flex-start;
        width: 100%;
    }

    .suggestion {
        display: flex;
        justify-content: flex-start;
        align-items: center;
        width: 100%;
        height: 64px;
        background-color: #ffffff;
        box-sizing: border-box;
        box-shadow: none;
        border: none;
        outline: none;
        border-top: 1px solid #aaaaaa;
        padding-left: 0.5em;
        transition: all 0.2s ease-in-out;

        &:hover, &:focus {
            cursor: pointer;
            background-color: #cccccc;
        }

        img {
            width: 48px;
            height: 48px;
        }

        .pharmacy-name {
            display: flex;
            text-wrap: nowrap;
            overflow: hidden;
            justify-content: center;
            align-items: center;
            flex: 1;

            & > * {
                font-size: 18px;
                user-select: none;
            }
        }
    }

    .search-box > input {
        flex: 1;
        font-size: large;
        height: 100%;
        outline: none;
        box-sizing: border-box;
        border: none;
    }
</style>