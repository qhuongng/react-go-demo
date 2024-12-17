import { Link } from "react-router-dom";
import { H1, P } from "../components/Typography";

const NotFound = () => {
    return (
        <div className="flex flex-col justify-center items-center min-h-screen">
            <H1>Whoops! You found uncharted territory!</H1>
            <P classNames="mt-24">How about we explore the area ahead of us later?</P>
            <Link to="/" className="btn btn-primary mt-8">
                Go back home
            </Link>
        </div>
    );
};

export default NotFound;
