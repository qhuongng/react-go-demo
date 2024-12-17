import { useState } from "react";
import { useForm, SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { RiEditLine } from "react-icons/ri";

import useAuthStore from "../../stores/auth.store";
import usePostStore from "../../stores/post.store";

import { H3, Sublabel } from "../Typography";
import LoadingIcon from "../LoadingIcon";

interface ModalProps {
    postId: number;
    postContent: string;
}

const PostEditSchema = z.object({
    content: z.string().nonempty({ message: "Type something first" }),
});

type PostEditInput = z.infer<typeof PostEditSchema>;

const PostEditModal: React.FC<ModalProps> = ({ postId, postContent }) => {
    const [open, setOpen] = useState(false);

    const accessToken = useAuthStore((state) => state.accessToken);
    const fetchPosts = usePostStore((state) => state.fetchPosts);

    const {
        register,
        handleSubmit,
        reset,
        formState: { errors, isSubmitting },
    } = useForm<PostEditInput>({
        defaultValues: { content: postContent },
        resolver: zodResolver(PostEditSchema),
    });

    const onSubmit: SubmitHandler<PostEditInput> = async (data) => {
        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL}/posts/${postId}`, {
                method: "PUT",
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
                const modalCheckbox = document.getElementById("edit-modal");
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
            <button className="btn btn-square btn-ghost" onClick={() => setOpen(true)}>
                <RiEditLine size={20} />
            </button>

            <input type="checkbox" id="edit-modal" className="modal-toggle" checked={open} />

            <div className="modal" role="dialog">
                <div className="modal-box">
                    <H3>Edit post</H3>

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
                                htmlFor="edit-modal"
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

export default PostEditModal;
