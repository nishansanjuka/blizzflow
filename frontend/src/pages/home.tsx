import { Button } from "@/components/ui/button";
import { useAuth } from "@/hooks/use-auth";

const Home: React.FC = () => {
  // const [name, setName] = useState<string>("");
  // const [time, setTime] = useState<string>("Listening for Time event...");
  const { logout } = useAuth();
  // const doGreet = () => {
  //   // let localName = name;
  //   // if (!localName) {
  //   //   localName = "anonymous";
  //   // }
  //   // GreetService.Greet(localName).then((resultValue: string) => {
  //   //   setResult(resultValue);
  //   // }).catch((err: any) => {
  //   //   console.log(err);
  //   // });
  // };

  // useEffect(() => {
  //   Events.On("time", (timeValue: any) => {
  //     setTime(timeValue.data);
  //   });
  //   // Reload WML so it picks up the wml tags
  //   WML.Reload();
  // }, []);

  return (
    <div>
      hey <Button onClick={logout}>Logout</Button>
    </div>
  );
};

export default Home;
