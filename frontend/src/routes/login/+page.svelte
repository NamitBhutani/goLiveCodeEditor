<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { enhance } from '$app/forms';
</script>

<Card.Root>
	<Card.Header>
		<Card.Title class="flex justify-center">Login</Card.Title>
	</Card.Header>
	<Card.Content>
		<div class="flex justify-center items-center mt-10">
			<form action="?/login" method="POST" use:enhance>
				<Label>Username:</Label>
				<Input type="text" name="username" id="username" />
				<Label>Password:</Label>
				<Input type="password" name="password" id="password" />
				<div class="flex justify-center mt-2"><Button>Login</Button></div>
			</form>
		</div>
	</Card.Content>
	<Card.Footer class="flex justify-center">
		<a href="/"><Button>Home</Button></a>
	</Card.Footer>
</Card.Root>
