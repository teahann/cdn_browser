<script lang=ts>
  import { onMount } from 'svelte'
  import { writable } from 'svelte/store'
  import { fly } from 'svelte/transition';
  import { createClient } from 'https://cdn.jsdelivr.net/npm/@supabase/supabase-js/+esm'

  const supabase = createClient(import.meta.env.VITE_PSQL_URL, import.meta.env.VITE_PSQL_KEY)
  const CDN = { raw: import.meta.env.VITE_CDN_RAW, short: import.meta.env.VITE_CDN_SHORT }
  const buckets = ['emoji_v2', 'stickers_v2']
  const library = writable(new Object())
  const gallery = writable(new Array())
  let selected_bucket = new String()
  const notification = {
    animation: { y: -100, opacity: 0, duration: 400 },
    duration: 2400,
    filename: new String(),
    show: false,
    timeout: null
  };

  onMount(async () => {
    try {
      const folders = await Promise.all(buckets.map(bucket => {
        return supabase.storage.from(bucket).list().then(({ data }) => {
          return { [bucket]: data?.map(item => item.name) || [] };
        });
      }));
      library.set(Object.assign({}, ...folders));
    } catch (err) { console.error("Database error") }
  });

  // Get all images from /selected_bucket/folder
  const get_images = async (folder) => {
    await supabase.storage.from(selected_bucket).list(folder).then(({ data, error }) => {
      gallery.set(data.map(item => `${CDN.raw}${selected_bucket}/${folder}/${item.name}`))
    })
  }

  // Copy image URL to clipboard, or open in new tab
  const select_img = async (cdn_url) => {
    let shareable_url = readable_url(cdn_url)
    notification.filename = img_name(shareable_url)
    try {
      await navigator.clipboard.writeText(shareable_url)
      copy_success()
    } catch (error) {
      window.open(shareable_url, '_blank')
    }
  }

  // Convert CDN direct URL into cleaner/shorter URL
  const readable_url = (cdn_url) => {
    let start_url = cdn_url.slice(66).split('/')[0]
    let end_url = cdn_url.slice(66).split('/').slice(1).join('/')
    return start_url === 'emoji_v2' ? `${CDN.short}emoji/${end_url}` : `${CDN.short}sticker/${end_url}`
  }

  // Handle "copied to clipboard" notification
  const copy_success = () => {
    if (notification.timeout) clearTimeout(notification.timeout)
    notification.show = true
    notification.timeout = setTimeout(() => {
      notification.show = false;
      notification.timeout = null;
    }, 2600);
  };

  // Replace underscores with spaces, capitalize words
  const clean_string = (str) => {
    return str.replace(/_/g, ' ')
      .split(' ').map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()).join(' ')
  }

  // Extract filename from URL 
  const img_name = (url) => url.substring(url.lastIndexOf('/') + 1)

</script>

{#if !$library.length && !selected_bucket.length}
  <main class="Home">
      {#each Object.keys($library) as bucket}
        <button on:click={() => selected_bucket = bucket}>{clean_string(bucket.slice(0,-3))}</button>
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
      <img src={image} on:click={() => select_img(image)} class={selected_bucket.slice(0,-3)} />
    {/each}
  </main>
{/if}

{#if notification.show}
  <div class="Modal" in:fly={notification.animation} out:fly={notification.animation}>
    <code class="Filename">{notification.filename}</code>
    <code class="Message">copied to clipboard!</code>
  </div>
{/if}
