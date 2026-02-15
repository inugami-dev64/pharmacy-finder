<script lang="ts">
    const {
        data
    } : {
        data: { slug: string}
    } = $props();
</script>

<h1>Hi there user: {data.slug}</h1>