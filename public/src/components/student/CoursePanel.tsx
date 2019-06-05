import * as React from "react";

import { DynamicTable } from "../../components";

import { IAssignment, IStudentSubmission } from "../../models";
import { Course } from "../../../proto/ag_pb";

import { NavigationManager } from "../../managers/NavigationManager";

interface IPanelProps {
    course: Course;
    labs: IStudentSubmission[];
    navMan: NavigationManager;
}
class CoursePanel extends React.Component<IPanelProps, any> {

    public render() {
        const labPath: string = "app/student/courses/" + this.props.course.getId() + "/lab/";
        const glabPath: string = "app/student/courses/" + this.props.course.getId() + "/grouplab/";

        return (
            <div className="col-lg-3 col-md-6 col-sm-6">
                <div className="panel panel-primary">
                    <div className="panel-heading clickable"
                        onClick={() => this.handleCourseClick()}>{this.props.course.getName()}</div>
                    <div className="panel-body">
                        <DynamicTable
                            header={["Labs", "Score", "Deadline"]}
                            data={this.props.labs}
                            selector={(item: IStudentSubmission) => {
                                const score = item.latest ? (item.latest.score.toString() + "%") : "N/A";
                                return [
                                    item.assignment.name,
                                    score,
                                    item.assignment.deadline.toDateString(),
                                ];
                            }}
                            onRowClick={(lab: IStudentSubmission) => {
                                this.handleRowClick(!lab.assignment.isgrouplab ? labPath : glabPath, lab.assignment);
                            }}
                        />
                    </div>
                </div>
            </div>
        );
    }

    private handleRowClick(pathPrefix: string, lab: IAssignment) {
        if (lab) {
            this.props.navMan.navigateTo(pathPrefix + lab.id);
        }
    }

    private handleCourseClick() {
        const uri: string = "app/student/courses/" + this.props.course.getId();
        this.props.navMan.navigateTo(uri);
    }
}

export { CoursePanel };
