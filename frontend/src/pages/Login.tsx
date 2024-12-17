import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import { useForm, SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { RiEyeLine, RiEyeOffLine } from "react-icons/ri";

import useAuthStore from "../stores/auth.store";
import { H2, Sublabel } from "../components/Typography";
import LoadingIcon from "../components/LoadingIcon";

const LoginSchema = z.object({
    username: z.string().nonempty({ message: "What's your username?" }),
    password: z
        .string()
        .min(6, { message: "Your password should be longer than 'TRAGIC' (6 characters)" }),
});

type LoginInput = z.infer<typeof LoginSchema>;

const Login = () => {
    const navigate = useNavigate();
    const setAuth = useAuthStore((state) => state.setAuth);

    const [showPassword, setShowPassword] = useState(false);

    const {
        register,
        handleSubmit,
        reset,
        formState: { errors, isSubmitting },
    } = useForm<LoginInput>({
        resolver: zodResolver(LoginSchema),
    });

    const onSubmit: SubmitHandler<LoginInput> = async (data) => {
        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
                body: JSON.stringify(data),
            });

            if (response.ok) {
                const data = await response.json();
                setAuth(data.data.id, data.data.accessToken);
                // login state to persist user login
                localStorage.setItem("isLoggedIn", "true");
                reset();
                navigate("/");
            } else {
                const data = await response.json();
                console.log(data.errors[0].message);
            }
        } catch (error) {
            console.log(error);
        }
    };

    return (
        <div className="flex flex-col justify-center items-center w-1/2 lg:w-1/3 py-12">
            <div className="flex flex-col justify-center items-center bg-base-200 border border-base-content rounded-box py-8 px-16 w-full">
                <H2>How's it goin'?</H2>

                <form className="mt-3 w-full" onSubmit={handleSubmit(onSubmit)}>
                    <label className="form-control w-full">
                        <div className="label">
                            <span className="label-text font-semibold">Username</span>
                            <span className="label-text font-semibold text-red-600">*</span>
                        </div>

                        <input
                            type="text"
                            placeholder="melon_usk"
                            className={`input w-full ${
                                errors.username ? "input-error" : "input-bordered"
                            }`}
                            {...register("username")}
                        />

                        <div className="label">
                            {errors.username && (
                                <Sublabel classNames="label-text-alt text-red-600">
                                    {errors.username.message}
                                </Sublabel>
                            )}
                        </div>
                    </label>

                    <label className="form-control w-full">
                        <div className="label">
                            <span className="label-text font-semibold">Password</span>
                            <span className="label-text font-semibold text-red-600">*</span>
                        </div>

                        <label
                            className={`input ${
                                errors.password ? "input-error" : "input-bordered"
                            } flex items-center gap-2`}
                        >
                            <input
                                type={showPassword ? "text" : "password"}
                                placeholder="••••••••"
                                className="grow"
                                {...register("password")}
                            />

                            <div className="btn btn-square btn-link -m-2">
                                {showPassword ? (
                                    <RiEyeLine size={20} onClick={() => setShowPassword(false)} />
                                ) : (
                                    <RiEyeOffLine size={20} onClick={() => setShowPassword(true)} />
                                )}
                            </div>
                        </label>

                        <div className="label">
                            {errors.password && (
                                <Sublabel classNames="label-text-alt text-red-600">
                                    {errors.password.message}
                                </Sublabel>
                            )}
                        </div>
                    </label>

                    <button
                        disabled={isSubmitting}
                        className="btn w-full btn-primary mt-6"
                        type="submit"
                    >
                        {isSubmitting ? <LoadingIcon /> : "Enter"}
                    </button>
                </form>
            </div>

            <Link className="mt-8 link link-hover" to="/signup">
                I'm new here
            </Link>

            <Link className="mt-4 link link-hover font-bold" to="/">
                JUST LET ME IN
            </Link>
        </div>
    );
};

export default Login;
