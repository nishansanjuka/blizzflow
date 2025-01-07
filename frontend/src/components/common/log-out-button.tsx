import { FC, memo, useCallback, useState } from "react";
import { Button } from "../ui/button";
import { Loader2, LogOut } from "lucide-react";
import { useAuth } from "@/hooks/use-auth";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "../ui/tooltip";

const LogOutButton: FC = () => {
  const [isLoading, setIsLoading] = useState(false);
  const { logout } = useAuth();
  const handleLogOut = useCallback(async () => {
    setIsLoading(true);
    await logout();
    setIsLoading(false);
  }, []);

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            disabled={isLoading}
            onClick={handleLogOut}
            variant={"ghost"}
            size={"icon"}
          >
            {isLoading ? <Loader2 className="animate-spin" /> : <LogOut />}
          </Button>
        </TooltipTrigger>
        <TooltipContent>
          <p>logout</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default memo(LogOutButton);
