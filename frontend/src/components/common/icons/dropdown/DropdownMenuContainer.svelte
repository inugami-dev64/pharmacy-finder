<script lang="ts">
    import { dropdownMenuZIndex } from "$lib/utils/z-indices";
    import type { Component, Snippet } from "svelte";

    let {
        size,
        ButtonComponent,
        title = undefined,
        color = "#999999",
        children = undefined
    }: {
        size: number,
        ButtonComponent: Component<{size: number, title: string, color?: string, padding?: number, onclick?: () => void, dummy?: boolean}>,
        title?: string,
        color?: string,
        children?: Snippet
    } = $props();

    let idx: string = Math.floor(Math.random() * Date.now()).toString(36);
</script>

<label for="_hidden-{idx}" class="dropdown-container" style="--size: {size}px">
    <input type="checkbox" name="_hidden" id="_hidden-{idx}" hidden>
    <ButtonComponent
        size={size}
        title={title || ""}
        padding={0}
        color={color}
        dummy={true}
    />

    <div class="dropdown" style="--zIndex: {dropdownMenuZIndex}">
        {@render children?.()}
    </div>
</label>

<style>
    .dropdown-container {
        display: inline-block;
        vertical-align: middle;
        padding: 0.5em;
        width: var(--size);
        height: var(--size);

        & > .dropdown {
            display: none;
            z-index: var(--zIndex);
        }

        @media (max-width: 600px) {
            & > input:checked ~ .dropdown {
                display: block;
                text-align: center;
                background-color: white;
                position: relative;
                border: 1px solid black;
                margin: 0;
                left: calc(var(--size) - 180px);
                width: 180px;
                border-radius: 1em;
                padding: 0;
                overflow: hidden;
                height: fit-content;
                box-sizing: border-box;

                :global(& > button), :global(& > a) {
                    all: unset;
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    user-select: none;
                    padding: 0;
                    height: 48px;
                    width: 100%;
                    transition: all 0.3s ease-in-out;

                    & > p {
                        flex: 1;
                    }

                    &:hover {
                        cursor: pointer;
                        background-color: rgba(0, 0, 0, 0.15)
                    }
                }
            }
        }

        &:hover > .dropdown {
            display: block;
            text-align: center;
            background-color: white;
            position: relative;
            border: 1px solid black;
            margin: 0;
            left: calc(var(--size) - 180px);
            width: 180px;
            border-radius: 1em;
            padding: 0;
            overflow: hidden;
            height: fit-content;
            box-sizing: border-box;

            :global(& > button), :global(& > a) {
                all: unset;
                display: flex;
                flex-direction: row;
                align-items: center;
                user-select: none;
                padding: 0;
                height: 48px;
                width: 100%;
                transition: all 0.3s ease-in-out;

                & > p {
                    flex: 1;
                }

                &:hover {
                    cursor: pointer;
                    background-color: rgba(0, 0, 0, 0.15)
                }
            }
        }
    }
</style>