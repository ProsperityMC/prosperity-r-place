<script lang="ts">
  import { Layer, t } from "svelte-canvas";

  const docImage = new Image();

  $: render = ({ context, width }) => {
    context.drawImage(docImage, 0, 0, 1024, 1024);
  };
</script>

<Layer {render} />
