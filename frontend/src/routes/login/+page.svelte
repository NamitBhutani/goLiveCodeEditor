<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
</script>

<Card.Root>
	<Card.Header>
		<Card.Title class="flex justify-center">Login</Card.Title>
	</Card.Header>
	<Card.Content>
		<form action="">
			<Label>Username:</Label>
			<Input type="text" />
			<Label>Password:</Label>
			<Input type="password" />
		</form>
	</Card.Content>
	<Card.Footer class="flex justify-center">
		<Button>Login</Button>
	</Card.Footer>
</Card.Root>
