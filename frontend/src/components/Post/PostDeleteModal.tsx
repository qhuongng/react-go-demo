import { RiDeleteBin6Line } from "react-icons/ri";

import useAuthStore from "../../stores/auth.store";
import usePostStore from "../../stores/post.store";

import { H3, P } from "../Typography";

interface ModalProps {
    postId: number;
}

const PostDeleteModal: React.FC<ModalProps> = ({ postId }) => {
    const accessToken = useAuthStore((state) => state.accessToken);
    const fetchPosts = usePostStore((state) => state.fetchPosts);

    const handleDelete = async () => {
        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL}/posts/${postId}`, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${accessToken}`,
                },
                credentials: "include",
            });

            if (response.ok) {
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
                className="btn btn-square btn-ghost -mr-2"
                onClick={() => {
                    const modal = document.getElementById("delete-modal");
                    if (modal) {
                        (modal as HTMLDialogElement).showModal();
                    }
                }}
            >
                <RiDeleteBin6Line size={20} />
            </div>

            <dialog id="delete-modal" className="modal">
                <div className="modal-box">
                    <H3>Woah!</H3>
                    <P>Are you sure you wanna get rid of this post? It's irreversible!</P>

                    <div className="modal-action">
                        <form method="dialog" className="flex gap-2">
                            {/* if there is a button in form, it will close the modal */}
                            <button className="btn btn-ghost">Naw</button>
                            <button className="btn btn-primary" onClick={handleDelete}>
                                Just do it
                            </button>
                        </form>
                    </div>
                </div>

                <form method="dialog" className="modal-backdrop">
                    {/* invisible button to enable closing the modal by clicking on the backdrop */}
                    <button></button>
                </form>
            </dialog>
        </div>
    );
};

export default PostDeleteModal;
