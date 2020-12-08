import * as React from "react";
import { Comment, User } from "../../../proto/ag_pb";


interface ICommentProps {
    comment?: Comment,
    author?: User,
    onSelect: () => void;
}

interface ICommentState {
    collapsed: boolean,
    editing: boolean,
}

export class IComment extends React.Component<ICommentProps, ICommentState>  {
    constructor(props: ICommentProps) {
        super(props);
        this.state = {
            collapsed: true,
            editing: false,
        }
    }

    public render() {
        if (!(this.props.comment && this.props.author)) {
            return <div></div>;
        }
        return <div className="row col-md-12 comment-body"
            onClick={() => this.props.onSelect()}
        >
            <div className="comment-pic col-md-3"><img src={this.props.author.getAvatarurl()} className="comment-img" /></div>

            <div className="row col-md-11 comment-header">
                <span className="coment-author">{this.props.author.getName()}</span>
                <span className="comment-posted">posted {this.props.comment.getPosted()}</span>
            </div>
            <div className="row col-md-11 comment-message">
                {this.props.comment.getMessage()}
            </div>
        </div>
    }
}