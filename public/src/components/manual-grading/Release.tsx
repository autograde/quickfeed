import * as React from "react";
import { Assignment, GradingBenchmark, GradingCriterion, Review, Submission, User } from '../../../proto/ag_pb';
import { totalScore, userSubmissionLink, submissionStatusToString, setDivider } from '../../componentHelper';
import { ISubmission } from "../../models";
import { formatDate } from '../../helper';
import ReactTooltip from "react-tooltip";

interface ReleaseProps {
    submission: ISubmission | undefined;
    assignment: Assignment;
    authorName: string;
    authorLogin: string;
    studentNumber: number;
    courseURL: string;
    teacherView: boolean;
    allClosed: boolean;
    setGrade: (status: Submission.Status, approved: boolean) => Promise<boolean>;
    release: (ready: boolean) => void;
    getReviewers: (submissionID: number) => Promise<User[]>;
    toggleCloseAll: () => void;
}

interface ReleaseState {
    open: boolean;
    reviews: Review[];
    reviewers: Map<User, Review>;
    score: number;
    status: Submission.Status;
}
export class Release extends React.Component<ReleaseProps, ReleaseState>{

    constructor(props: ReleaseProps) {
        super(props);
        this.state = {
            reviews: [],
            score: totalScore(this.selectReadyReviews()),
            reviewers: new Map<User, Review>(),
            open: !this.props.teacherView,
            status: Submission.Status.NONE,
        }
    }

    public render() {
        const open = this.state.open && !this.props.allClosed;
        const reviewInfoSpan = <span className="r-info">Reviews: {this.props.submission?.reviews.length ?? 0}/{this.props.assignment.getReviewers()}</span>;
        const noReviewsSpan = <span className="r-info">N/A</span>;
        const noSubmissionDiv = <div className="alert alert-info">No submissions for {this.props.assignment.getName()}</div>;
        const noReviewsDiv = <div className="alert alert-info">{this.props.assignment.getName()} is not for manual grading</div>
        const noReadyReviewsDiv = <div className="alert alert-info">No ready reviews for {this.props.assignment.getName()}</div>

        const headerDiv = <div className="row review-header" onClick={() => this.toggleOpen()}>
        <h3><span className="r-number">{this.props.studentNumber}. </span><span className="r-header">{this.props.authorName}</span><span className="r-score">Score: {totalScore(this.props.submission?.reviews ?? [])} </span>{this.props.assignment.getReviewers() > 0 ? reviewInfoSpan : noReviewsSpan}{this.releaseButton()}</h3>
        </div>;

        if (this.props.assignment.getReviewers() < 1) {
            return <div className="release">
                {headerDiv}
                {open ? noReviewsDiv : null}
            </div>;
        }

        if (!this.props.submission) {
            return <div className="release">
                {headerDiv}
                {open ? noSubmissionDiv : null}
            </div>;
        }

        if (this.state.reviews.length < 1) {
            return <div className="release">
                {headerDiv}
                {open ? noReadyReviewsDiv : null}
            </div>;
        }

        return <div className="release">
            {this.props.teacherView ? headerDiv : null}
            {open ? setDivider() : null}
            {open && this.props.teacherView ? this.infoTable() : null}
            {open ? this.renderReleaseTable() : null}
            {open}
        ></div>;
    }

    private infoTable(): JSX.Element {
        return <div className="row">
            <div className="col-md-6 release-info">
                <ul className="list-group">
                    <li key="li0" className="list-group-item r-li">
                        <span className="r-table">Deadline: </span>
                            {formatDate(this.props.assignment.getDeadline())}</li>
                    <li key="li1" className="list-group-item r-li">
                        <span className="r-table">Delivered: </span>
                            {this.props.submission ? formatDate(this.props.submission?.buildDate) : "Not delivered"}</li>
                    <li key="li3" className="list-group-item r-li">
                        <span className="r-table">Repository: </span>
                        {userSubmissionLink(this.props.authorLogin, this.props.assignment.getName(), this.props.courseURL, "btn btn-default")}</li>
                    <li key="li4" className="list-group-item r-li"><span className="r-table">Status: </span>{submissionStatusToString(this.state.status)}</li>
                    <li key="li5" className="list-group-item r-li">{this.renderStatusButton()}</li>
                </ul>
            </div>
            <div className="col-md-6">
                <table className="table">
                    <thead><tr key="it">
                            <td>Reviewers:</td>
                            <td>Score:</td>
                        </tr></thead>
                        <tbody>
                        {Array.from(this.state.reviewers.keys()).map((r, i) => <tr key={"it" + i}>
                            <td>{r.getName()}</td>
                            <td>{this.state.reviewers.get(r)?.getScore() ?? 0}</td>
                        </tr>)}</tbody>
                </table>
            </div>
        </div>;
    }

    private releaseButton(): JSX.Element {
        return <div
            className={this.releaseButtonClass()}
            onClick={() => {
                if (this.props.submission && this.props.assignment.getReviewers() > 0) {
                    this.props.release(!this.props.submission.released);
                }
            }}
            >{this.releaseButtonString()}</div>;
        }

    private releaseButtonClass(): string {
        if (!this.props.submission || this.props.assignment.getReviewers() < 1 ||
         this.props.submission.reviews.length < this.props.assignment.getReviewers()) {
             return "btn release-btn";
         }
        return "btn btn-default release-btn";
    }

    private releaseButtonString(): string {
        if (!this.props.submission || this.props.assignment.getReviewers() < 1) {
             return "N/A";
         }
        return this.props.submission.released ? "Released" : "Release";
    }

    private selectReadyReviews(): Review[] {
        const selected: Review[] = [];
        this.props.submission?.reviews.forEach(r => {
            if (r.getReady()) selected.push(r);
        });
        return selected;
    }

    private renderReleaseTable(): JSX.Element {
        const reviewersList = Array.from(this.state.reviewers.keys());
        return <div className="row">
            <table className="table table-condensed table-bordered">
            <thead><tr key="rthead"><th>Reviews:</th>{reviewersList.map((u, i) => <th className="release-cell">
                {(this.state.reviewers.get(u)?.getScore() ?? 0) + "%"}
            </th>)}</tr></thead>
            <tbody>
                {this.renderTableRows()}
            </tbody>
            </table>
        </div>;
    }

    private renderTableRows(): JSX.Element[] {
        const rows: JSX.Element[] = [];
        const reviewersList = Array.from(this.state.reviewers.keys());
        this.props.assignment.getGradingbenchmarksList().forEach((bm, i) => {
            rows.push(<tr key={"rt" + i} className="b-header"><td>{bm.getHeading()}</td>{reviewersList.map(u =>
                <td>{this.commentSpan(this.selectBenchmark(u, bm).getComment(), "bm" + bm.getId())}</td>)}</tr>);
            bm.getCriteriaList().forEach((c, j) => {
                rows.push(<tr key={"rrt" + j + i}><td>{c.getDescription()}</td>
                {reviewersList.map(u => <td className={this.setCellColor(u, c)}>
                    <span className={this.setCellIcon(u, c)}></span>
                    {this.commentSpan(this.selectCriterion(u, c).getComment(), "cr" + c.getId())}
                </td>)}
                </tr>);
            });
        });
        rows.push(<tr key="rtf"><td>Feedbacks:</td>
            {reviewersList.map((u, i) => <td>{this.commentSpan(this.state.reviewers.get(u)?.getFeedback() ?? "No feedback", "fb" + i)}</td>)}
        </tr>);
        rows.push(<tr key="tscore"><td>Score: {this.props.submission?.score ?? 0}</td>
            {reviewersList.map(u => <td>{this.state.reviewers.get(u)?.getScore() ?? 0}</td>)}
        </tr>);
        return rows;
    }

    private setCellIcon(u: User, c: GradingCriterion): string {
        const cr = this.selectCriterion(u, c);
        switch (cr.getGrade()) {
            case GradingCriterion.Grade.PASSED:
                return "r-cell glyphicon glyphicon-ok";
            case GradingCriterion.Grade.FAILED:
                return "r-cell glyphicon glyphicon-remove";
            default:
                return "r-cell glyphicon glyphicon-ban-circle";
        }
    }

    private setCellColor(u: User, c: GradingCriterion): string {
            const cr = this.selectCriterion(u, c);
            if (cr.getGrade() === GradingCriterion.Grade.PASSED) {
                return "success";
            }
            return cr.getGrade() === GradingCriterion.Grade.FAILED ? "danger" : "";
    }

    private selectBenchmark(u: User, bm: GradingBenchmark): GradingBenchmark {
        const r = this.state.reviewers.get(u);
        if (r) {
            const rbm = r.getBenchmarksList().find(item => item.getId() === bm.getId())
            if (rbm) bm = rbm;
        }
        return bm;
    }

    private selectCriterion(u: User, c: GradingCriterion): GradingCriterion {
        const r = this.state.reviewers.get(u);
        if (r) {
            r.getBenchmarksList().forEach(bm => {
                const rc = bm.getCriteriaList().find(item => item.getId() === c.getId());
                if (rc) c = rc;
            });
        }
        return c;
    }

    private commentSpan(text: string, id: string): JSX.Element {
        if (text === "") {
            return <span></span>;
        }
        return <span><span className="release-comment glyphicon glyphicon-comment"
            data-tip
            data-for={id}
        ></span>
        <ReactTooltip
            type="light"
            effect="solid"
            id={id}
        ><p>{text}</p></ReactTooltip></span>;
    }

    private renderStatusButton(): JSX.Element {
        return <div className="form-group r-grade">
            <select className="form-control" onChange={(e) => this.updateStatus(e.target.value)}>
                <option key="st0" value="none" selected={this.props.submission?.status === Submission.Status.NONE}>Set status</option>
                <option key="st1" value="approve" selected={this.props.submission?.status === Submission.Status.APPROVED}>Approved</option>
                <option key="st2" value="reject" selected={this.props.submission?.status === Submission.Status.REJECTED}>Rejected</option>
                <option key="st3" value="revision" selected={this.props.submission?.status === Submission.Status.REVISION}>Revision</option>
            </select>
            </div>;
    }

    private async updateStatus(action: string) {
        if (this.props.submission) {
            let newStatus: Submission.Status = Submission.Status.NONE;
            let newBool = false;
            switch (action) {
                case "approve":
                    newStatus = Submission.Status.APPROVED;
                    newBool = true;
                    break;
                case "reject":
                    newStatus = Submission.Status.REJECTED;
                    break;
                case "revision":
                    newStatus = Submission.Status.REVISION;
                    break;
                default:
                    newStatus = Submission.Status.NONE;
                    break;
            }
            const ans = this.props.setGrade(newStatus, newBool);
            if (ans) {
                this.setState({
                    status: newStatus,
                })
            }

        }
    }

    private async mapReviewers(): Promise<Map<User, Review>> {
        const reviews = this.selectReadyReviews();
        const updatedMap = new Map<User, Review>();
        if (this.props.submission && reviews.length > 0) {
            const reviewers = await this.props.getReviewers(this.props.submission.id);
            reviewers.forEach(r => {
                const selectedReview = this.selectReviewByReviewer(r, reviews);
                if (selectedReview) updatedMap.set(r, selectedReview);
            });
        }
        return updatedMap;
    }

    private selectReviewByReviewer(user: User, reviews: Review[]): Review | undefined {
        return reviews.find(item => item.getReviewerid() === user.getId());
    }

    private async toggleOpen() {
        // if closing, flush the state
        if (this.state.open) {
            this.setState({
                reviews: [],
                reviewers: new Map<User, Review>(),
                open: false,
            })
        }

        // if opening, close all other reviews
        if (!(this.state.open && this.props.allClosed)) {
            this.props.toggleCloseAll();
        }
        const ready = this.selectReadyReviews();
        if (ready.length > 0) {
            this.setState({
                open: !this.state.open,
                reviewers: await this.mapReviewers(),
                reviews: ready,
                score: totalScore(ready),
                status: this.props.submission?.status ?? Submission.Status.NONE,
            });
        } else {
            this.setState({open: !this.state.open});
        }
    }
}