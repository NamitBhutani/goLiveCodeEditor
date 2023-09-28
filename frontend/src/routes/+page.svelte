<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
</script>

<Card.Root>
	<Card.Header>
		<Card.Title class="flex justify-center">goLiveCodeEditor</Card.Title>
	</Card.Header>
	<Card.Content>
		<div class="flex justify-center">Edit away!</div>
	</Card.Content>
	<Card.Footer class="flex justify-center">
		<a href="/login"><Button class="mx-2">Login</Button></a>
		<a href="/register"><Button class="mx-2">Register</Button></a>
	</Card.Footer>
</Card.Root>
