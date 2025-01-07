import { Loader2, Snowflake } from "lucide-react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAuth } from "@/hooks/use-auth";
import { toast } from "sonner";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { Input } from "@/components/ui/input";
import { SetWIndowFullScreen } from "@/lib/utils";

const loginSchema = z.object({
  username: z.string().min(3, "Username must be at least 3 characters"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

type LoginFormData = z.infer<typeof loginSchema>;

const Login: React.FC = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);

  const { login } = useAuth();

  const onSubmit = async (data: LoginFormData) => {
    try {
      setLoading(true);
      toast.loading("Logging in...");
      await new Promise((resolve) => setTimeout(resolve, 2000));
      await login(data.username, data.password);
      toast.success("Logged in successfully");
      SetWIndowFullScreen();
      navigate("/");
    } catch (error) {
      if (typeof error === "string") {
        toast.error(error.split(":")[1]);
      }
    } finally {
      setLoading(false);
      toast.dismiss();
    }
  };

  return (
    <div className="fixed w-full h-screen flex items-center justify-center">
      <div className="p-8 rounded w-[90%]">
        <div className="w-full flex flex-col justify-center items-center mb-4">
          <Snowflake className="size-10 text-blue-500" />
          <span className="text-xl font-extrabold uppercase text-blue-500 mb-10">
            Blizzflow
          </span>
          <h1 className="text-2xl font-bold text-center">Welcome Back!</h1>
          <p className="text-gray-600 text-center mt-2">
            Sign in to manage your workflow seamlessly
          </p>
        </div>
        <form className="mt-4" onSubmit={handleSubmit(onSubmit)}>
          <div className="mb-4">
            <label
              htmlFor="username"
              className="block text-sm font-medium text-gray-700"
            >
              Username
            </label>
            <Input {...register("username")} type="text" id="username" />
            {errors.username && (
              <p className="mt-1 text-sm text-red-600">
                {errors.username.message}
              </p>
            )}
          </div>
          <div className="mb-4">
            <label
              htmlFor="password"
              className="block text-sm font-medium text-gray-700"
            >
              Password
            </label>
            <Input {...register("password")} type="password" id="password" />
            {errors.password && (
              <p className="mt-1 text-sm text-red-600">
                {errors.password.message}
              </p>
            )}
          </div>
          <Button
            disabled={loading}
            type="submit"
            variant={"default"}
            className="w-full bg-blue-500 hover:bg-blue-400"
          >
            {loading ? <Loader2 className="size-4 animate-spin" /> : "Login"}
          </Button>
        </form>
      </div>
    </div>
  );
};

export default Login;
