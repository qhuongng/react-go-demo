interface HeadingsProps {
    classNames?: string;
    children: React.ReactNode;
}

export const H1: React.FC<HeadingsProps> = ({ classNames, children }) => {
    return (
        <h1
            className={`mb-4 text-3xl font-bold leading-none tracking-tight lg:text-4xl ${classNames}`}
        >
            {children}
        </h1>
    );
};

export const H2: React.FC<HeadingsProps> = ({ classNames, children }) => {
    return (
        <h2
            className={`mb-4 text-2xl font-semibold leading-none tracking-tight lg:text-3xl ${classNames}`}
        >
            {children}
        </h2>
    );
};

export const H3: React.FC<HeadingsProps> = ({ classNames, children }) => {
    return (
        <h3
            className={`mb-4 text-xl font-medium leading-none tracking-tight lg:text-2xl ${classNames}`}
        >
            {children}
        </h3>
    );
};

export const H4: React.FC<HeadingsProps> = ({ classNames, children }) => {
    return (
        <h4
            className={`mb-3 text-md font-medium leading-none tracking-tight lg:text-lg ${classNames}`}
        >
            {children}
        </h4>
    );
};

export const P: React.FC<HeadingsProps> = ({ classNames, children }) => {
    return <p className={`text-md tracking-tight lg:text-lg ${classNames}`}>{children}</p>;
};

export const Sublabel: React.FC<HeadingsProps> = ({ classNames, children }) => {
    return <h4 className={`text-xs tracking-tight ${classNames}`}>{children}</h4>;
};
