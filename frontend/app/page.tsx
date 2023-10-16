import {block} from "million/react"
import {Button} from "@/components/ui/button"


const Block = block( function Home() {
  return (
    <div>
      <h1>Home</h1>
      <Button>Button</Button>
    </div>
  )
})

export default function App() {
  return (
    <div>
      <Block />
    </div>
  );
}
