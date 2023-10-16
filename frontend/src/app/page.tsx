import { Metadata } from 'next'
import { block } from "million/react-server";
export const metadata: Metadata = {
  title: 'My Page Title',
}
const Main = () => {
  
  return (
    <main >
      <h1 >Welcome to Next.js with Million!</h1>

      <p >
        Get started by editing{" "}
        <code >pages/index.tsx</code>
      </p>

      <p >
        Check out <a href="millionjs.org">Millionjs</a>
      </p>

    
    </main>
  );
};

const MainBlock = block(Main);

const Home = () => {
  return (
    <div >
      <MainBlock />
    </div>
  );
};

export default Home;