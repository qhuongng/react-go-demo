import useAuthStore from "../../stores/auth.store";

import PostEditModal from "./PostEditModal";
import PostDeleteModal from "./PostDeleteModal";
import { H3, P, Sublabel } from "../Typography";

import { Post } from "../../lib/types";
import { formatDateString } from "../../lib/utils";

interface PostProps {
    post: Post;
}

const PostContainer: React.FC<PostProps> = ({ post }) => {
    const userId = useAuthStore((state) => state.id);

    return (
        <div className="flex flex-col justify-center items-start bg-gray-100 border border-gray-600 rounded-box py-6 px-8 w-full">
            <div className="flex justify-between w-full">
                <H3>{post.userName}</H3>
                {userId === post.userId && (
                    <div className="flex gap-2">
                        <PostEditModal key={post.id} postId={post.id} postContent={post.content} />
                        <PostDeleteModal key={post.id} postId={post.id} />
                    </div>
                )}
            </div>

            <P classNames="mb-2">{post.content}</P>

            <Sublabel classNames="self-end">
                Posted on: {formatDateString(post.createdAt)}{" "}
                {post.updatedAt !== post.createdAt
                    ? `(edited on: ${formatDateString(post.updatedAt)})`
                    : ""}
            </Sublabel>
        </div>
    );
};

export default PostContainer;
