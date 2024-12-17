import { useState } from "react";
import { useForm, SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { RiAddLargeLine } from "react-icons/ri";

import useAuthStore from "../../stores/auth.store";
import usePostStore from "../../stores/post.store";

import { H3, Sublabel } from "../Typography";
import LoadingIcon from "../LoadingIcon";

const PostEditSchema = z.object({
    content: z.string().nonempty({ message: "Type something first" }),
});

type PostEditInput = z.infer<typeof PostEditSchema>;

const PostCreateModal: React.FC = () => {
    const [open, setOpen] = useState(false);

    const accessToken = useAuthStore((state) => state.accessToken);
    const fetchPosts = usePostStore((state) => state.fetchPosts);

    const {
        register,
        handleSubmit,
        reset,
        formState: { errors, isSubmitting },
    } = useForm<PostEditInput>({
        resolver: zodResolver(PostEditSchema),
    });

    const onSubmit: SubmitHandler<PostEditInput> = async (data) => {
        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL}/posts`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${accessToken}`,
                },
                credentials: "include",
                body: JSON.stringify(data),
            });

            if (response.ok) {
                setOpen(false);
                reset();
                const modalCheckbox = document.getElementById("add-modal");
                if (modalCheckbox) {
                    (modalCheckbox as HTMLInputElement).checked = false;
                }
                fetchPosts();
            } else {
                const data = await response.json();
                console.log(data.errors[0].message);
            }
        } catch (error) {
            console.log(error);
        }
    };

    return (
        <div>
            <div
                className="tooltip tooltip-left z-[50] fixed bottom-8 right-8"
                data-tip="Make spam"
            >
                <button
                    className="btn btn-circle btn-primary lg:btn-lg"
                    onClick={() => setOpen(true)}
                >
                    <RiAddLargeLine size={20} />
                </button>
            </div>

            <input type="checkbox" id="add-modal" className="modal-toggle" checked={open} />

            <div className="modal" role="dialog">
                <div className="modal-box">
                    <H3>Create a post</H3>

                    <form onSubmit={handleSubmit(onSubmit)}>
                        <textarea
                            className={`textarea ${
                                errors.content ? "textarea-error" : "textarea-bordered"
                            } w-full`}
                            placeholder="The meaning of life is 42"
                            {...register("content")}
                        />
                        {errors.content && (
                            <Sublabel classNames="text-red-600">{errors.content.message}</Sublabel>
                        )}

                        <div className="modal-action">
                            <label
                                htmlFor="add-modal"
                                className="btn btn-ghost"
                                onClick={() => {
                                    reset();
                                    setOpen(false);
                                }}
                            >
                                Scrap it
                            </label>

                            <button
                                className="btn btn-primary"
                                disabled={isSubmitting}
                                type="submit"
                            >
                                {isSubmitting ? <LoadingIcon /> : "I'm done"}
                            </button>
                        </div>
                    </form>
                </div>

                <button
                    className="modal-backdrop"
                    onClick={() => {
                        reset();
                        setOpen(false);
                    }}
                />
            </div>
        </div>
    );
};

export default PostCreateModal;
