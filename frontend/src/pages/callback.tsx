import { Loader } from "lucide-react";

function CallbackPage() {
  return (
    <div className="w-full h-screen flex text-muted-foreground items-center justify-center flex-col animate-pulse">
      <Loader className=" animate-spin" />
      <span>please wait...</span>
    </div>
  );
}

export default CallbackPage;
