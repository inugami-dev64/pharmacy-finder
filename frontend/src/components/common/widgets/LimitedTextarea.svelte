<script lang="ts">
    let {
        name,
        placeholder,
        text,
        maxLength = 1024
    }: {
        name: string,
        placeholder: string,
        text: string | undefined,
        maxLength: number | undefined
    } = $props();

    let currentInputLength = $state(text?.length || 0)
</script>

<div class="limited-textarea-container">
    <textarea
        bind:value={text}
        oninput={_ => currentInputLength = text?.length || 0}
        maxlength="{maxLength}"
        name="{name}"
        placeholder="{placeholder}"
    ></textarea>
    <p>{currentInputLength}/{maxLength}</p>
</div>

<style>
    .limited-textarea-container {
        width: 100%;
        border-radius: 5px;
        border: 1px solid #aaa;
        box-sizing: border-box;
        transition: all 0.2s;

        &:focus-within {
            border-width: 2px;
            border-color: #00c;
        }

        & > textarea {
            all: unset;
            margin: 0;
            padding: 3px;
            resize: none;
            width: 100%;
            min-height: 150px;
            box-sizing: border-box;
        }

        & > p {
            user-select: none;
            padding: 0;
            margin: 0;
            margin-right: 5px;
            text-align: right;
        }
    }
</style>