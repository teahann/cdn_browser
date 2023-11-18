<script>
  import { onMount } from 'svelte'
  import { writable } from 'svelte/store'
  import { createClient } from 'https://cdn.jsdelivr.net/npm/@supabase/supabase-js/+esm'

  const db = createClient(import.meta.env.VITE_PSQL_URL, import.meta.env.VITE_PSQL_KEY)
  const CDN = import.meta.env.VITE_CDN
  const buckets = ['emoji', 'stickers']
  const library = writable(new Array())
  let selected_bucket = new String()
  const gallery = writable(new Array())

  onMount(async () => {
    const results = await Promise.all(buckets.map(bucket => db.storage.from(bucket).list()));
    library.set(buckets.reduce((acc, bucket, i) => { 
      acc[bucket] = results[i].data.map(item => item.name)
      return acc
    }, {}))
  });

  const get_images = async (category) => {
    const { data } = await db.storage.from(selected_bucket).list(category).then(({ data, error }) => {
      const names = data.map(item => `${CDN}${selected_bucket}/${category}/${item.name}`)
      gallery.set(names)
    })
  }

  const hdl_image = async (url) => {
    try {
      await navigator.clipboard.writeText(url)
    } catch (error) {
      window.open(url, '_blank')
    }
  }

  const clean_string = (str) => {
    return str
      .replace(/_/g, ' ')
      .split(' ')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(' ')
  }

</script>

{#if !$library.length && !selected_bucket.length}
  <main class="Home">
      {#each Object.keys($library) as bucket}
        <button on:click={() => selected_bucket = bucket}>{clean_string(bucket)}</button>
      {/each}
  </main>
{:else if !$gallery.length}
  <main class="List">
    <button id="Back_Button" on:click={() =>  selected_bucket = new String()}>Back</button>
    {#each $library[selected_bucket] as category}
      {#if category !== '.emptyFolderPlaceholder'}
      <button on:click={() => get_images(category)}>{clean_string(category)}</button>
      {/if}
    {/each}
  </main>
{:else}
  <main class="List Gallery">
    <button id="Back_Button" on:click={() => gallery.set(new Array())}>Back</button>
    {#each $gallery as image}
      <img src={image} on:click={() => hdl_image(image)} class={selected_bucket} />
    {/each}
  </main>
{/if}
