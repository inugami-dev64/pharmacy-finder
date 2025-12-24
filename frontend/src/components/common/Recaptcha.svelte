<script lang="ts">
    import { onMount } from "svelte";

    const siteKey = process.env.RECAPTCHA_SITE_KEY;

    onMount(() => {
        const s = document.createElement("script");
        s.type = "text/javascript";
        s.async = true;
        s.defer = true;
        s.src = "https://www.google.com/recaptcha/api.js";
        s.onerror = () => {
            console.error("Recaptcha error");
        }

        if (siteKey) document.head.appendChild(s);

        return () => {
            s.remove();
        }
    });
</script>

<div class="g-recaptcha" data-sitekey="{siteKey}"></div>