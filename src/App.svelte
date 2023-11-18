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
      return acc;
    }, {}))
  });

  const get_images = async (category) => {
    const { data } = await db.storage.from(selected_bucket).list(category).then(({ data, error }) => {
      const names = data.map(item => make_url(category, item.name));
      gallery.set(names)
    })
  }

  const view_bucket = (bucket) => selected_bucket = bucket
  const make_url = (category, filename) => `${CDN}${selected_bucket}/${category}/${filename}`

  const exit_bucket = () => selected_bucket = new String()
  const exit_gallery = () => gallery.set(new Array())

  const hdl_image = async (url) => {
    try {
      await navigator.clipboard.writeText(url)
    } catch (error) {
      window.open(url, '_blank')
    }
  }

</script>

{#if !$library.length && !selected_bucket.length}
  <div class="Home">
      {#each Object.keys($library) as bucket}
        <button on:click={() => view_bucket(bucket)}>
          {bucket}
        </button>
      {/each}
  </div>
{:else if !$gallery.length}
  <div class="Bucket">
    <button on:click={() => exit_bucket()}>Back</button>
    {#each $library[selected_bucket] as category}
      <button on:click={() => get_images(category)}>{category}</button>
    {/each}
  </div>
{:else}
  <div class="Gallery">
    <button on:click={() => exit_gallery()}>Back</button>
    {#each $gallery as image}
      <img src={image} on:click={() => hdl_image(image)}/>
    {/each}
  </div>
{/if}
