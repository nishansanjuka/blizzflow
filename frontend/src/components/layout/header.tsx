import { FC } from "react";
import LogOutButton from "../common/log-out-button";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "../ui/breadcrumb";
import { SidebarTrigger } from "../ui/sidebar";
import { Separator } from "../ui/separator";
import { useAuth } from "@/hooks/use-auth";
export const Header: FC = () => {
  const { user, isAuthenticated } = useAuth();
  if (isAuthenticated && user) {
    return (
      <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
        <SidebarTrigger className="-ml-1" />
        <Separator orientation="vertical" className="mr-2 h-4" />
        <Breadcrumb className="flex-1">
          <BreadcrumbList>
            <BreadcrumbItem className="hidden md:block">
              <BreadcrumbLink href="#">Blizzflow</BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator className="hidden md:block" />
            <BreadcrumbItem>
              <BreadcrumbPage>Point Sales System</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
        <span className="flex items-center space-x-2">
          <span className="text-green-100 py-1 px-4 bg-green-500/70 rounded-md">
            {user.Username || "blizzflow"}
          </span>
          <LogOutButton />
        </span>
      </header>
    );
  }
};
