import * as React from "react";
import { CoursesOverview, NavMenu, SingleCourseOverview, StudentLab } from "../components";

import { CourseManager } from "../managers/CourseManager";
import { ILink, NavigationManager } from "../managers/NavigationManager";
import { UserManager } from "../managers/UserManager";

import {
    CourseUserState,
    IAssignment,
    ICourse,
    ICoursesWithAssignments,
    ICourseUserLink,
    IStudentSubmission,
    IUser,
    IUserCourse,
} from "../models";

import { View, ViewPage } from "./ViewPage";
import { HelloView } from "./views/HelloView";
import { UserView } from "./views/UserView";

import { ArrayHelper } from "../helper";
import { INavInfo, INavInfoEvent } from "../NavigationHelper";

import { CollapsableNavMenu } from "../components/navigation/CollapsableNavMenu";
import { ILinkCollection } from "../managers";
import { EnrollmentView } from "./views/EnrollmentView";

export class StudentPage extends ViewPage {
    private navMan: NavigationManager;
    private userMan: UserManager;
    private courseMan: CourseManager;

    private courses: IUserCourse[] = [];
    private activeCourses: IUserCourse[] = [];
    private selectedCourse: IUserCourse | undefined;
    private selectedAssignment: IStudentSubmission | undefined;

    private foundId: number = -1;

    constructor(users: UserManager, navMan: NavigationManager, courseMan: CourseManager) {
        super();

        this.navMan = navMan;
        this.userMan = users;
        this.courseMan = courseMan;

        this.navHelper.defaultPage = "index";

        this.navHelper.registerFunction<any>("index", this.index);
        this.navHelper.registerFunction<any>("course/{courseid}", this.course);
        this.navHelper.registerFunction<any>("course/{courseid}/lab/{labid}", this.courseWithLab);
        this.navHelper.registerFunction<any>("course/{courseid}/{page}", this.courseMissing);
        this.navHelper.registerFunction<any>("enroll", this.enroll);

        // Only for testing purposes
        this.navHelper.registerFunction<any>("user", this.getUsers);
        this.navHelper.registerFunction<any>("hello", (navInfo) => Promise.resolve(<HelloView></HelloView>));
    }

    public async getUsers(navInfo: INavInfo<any>): View {
        await this.setupData();
        return <UserView users={await this.userMan.getAllUser()}>
        </UserView>;
    }

    public async index(navInfo: INavInfo<any>): View {
        await this.setupData();
        if (this.activeCourses) {
            console.log(this.activeCourses);
            return (<CoursesOverview
                courseOverview={this.activeCourses}
                navMan={this.navMan}
            />);
        }
        return <h1>404</h1>;
    }

    public async enroll(navInfo: INavInfo<any>): View {
        await this.setupData();
        const curUser = this.userMan.getCurrentUser();
        if (!curUser) {
            return <h1>404</h1>;
        }
        return <div>
            <h1>Enrollment page</h1>
            <EnrollmentView
                courses={this.courses}
                onEnrollmentClick={(course: ICourse) => {
                    this.courseMan.addUserToCourse(curUser, course);
                    this.navMan.refresh();
                }}>
            </EnrollmentView>
        </div >;
    }

    public async course(navInfo: INavInfo<{ courseid: string }>): View {
        await this.setupData();
        this.selectCourse(navInfo.params.courseid);
        if (this.selectedCourse) {
            return (<SingleCourseOverview courseAndLabs={this.selectedCourse} />);
        }
        return <h1>404 not found</h1>;
    }

    public async courseWithLab(navInfo: INavInfo<{ courseid: string, labid: string }>): View {
        await this.setupData();
        this.selectCourse(navInfo.params.courseid);
        console.log("Course with lab", this.selectedCourse);
        if (this.selectedCourse) {
            await this.selectAssignment(navInfo.params.labid);
            if (this.selectedAssignment) {
                console.log("selected!");
                return <StudentLab
                    course={this.selectedCourse.course}
                    assignment={this.selectedAssignment}>
                </StudentLab>;
            }
        }
        console.log(navInfo);
        return <div>404 not found</div>;
    }

    public async courseMissing(navInfo: INavInfo<{ courseid: string, page: string }>): View {
        return <div>The page {navInfo.params.page} is not yet implemented</div >;
    }

    public async renderMenu(key: number): Promise<JSX.Element[]> {
        if (key === 0) {
            const coursesLinks: ILinkCollection[] = this.activeCourses.map(
                (course, i) => {
                    const allLinks: ILink[] = [];
                    allLinks.push({ name: "Labs" });
                    const labs = course.assignments;
                    allLinks.push(...labs.map((lab, ind) => {
                        return {
                            name: lab.assignment.name,
                            uri: this.pagePath + "/course/" + course.course.id + "/lab/" + lab.assignment.id,
                        };
                    }));
                    allLinks.push({ name: "Group Labs" });
                    allLinks.push({ name: "Settings" });
                    allLinks.push({
                        name: "Members", uri: this.pagePath + "/course/" + course.course.id + "/members",
                    });
                    allLinks.push({
                        name: "Coruse Info", uri: this.pagePath + "/course/" + course.course.id + "/info",
                    });
                    return {
                        item: { name: course.course.code, uri: this.pagePath + "/course/" + course.course.id },
                        children: allLinks,
                    };
                });

            const settings = [
                { name: "Join course", uri: this.pagePath + "/enroll" },
            ];

            this.navMan.checkLinkCollection(coursesLinks, this);
            this.navMan.checkLinks(settings, this);

            return [
                <h4 key={0}>Courses</h4>,
                <CollapsableNavMenu key={1} links={coursesLinks} onClick={(link) => this.handleClick(link)}>
                </CollapsableNavMenu>,
                <h4 key={2}>Settings</h4>,
                <NavMenu key={3} links={settings} onClick={(link) => this.handleClick(link)}></NavMenu>,
            ];
        }
        return [];
    }

    private onlyActiveCourses(studentCourse: IUserCourse[]): IUserCourse[] {
        const temp: IUserCourse[] = [];
        studentCourse.forEach((a) => {
            if (a.link && a.link.state === CourseUserState.student) {
                temp.push(a);
            }
        });
        return temp;
    }

    private async setupData() {
        const curUser = this.userMan.getCurrentUser();
        console.log("Setup data");
        if (curUser) {
            this.courses = await this.courseMan.getStudentCourses(curUser);
            this.activeCourses = this.onlyActiveCourses(this.courses);
        }
    }

    private selectCourse(courseId: string) {
        this.selectedCourse = undefined;
        const course = parseInt(courseId, 10);
        if (!isNaN(course)) {
            this.selectedCourse = this.activeCourses.find(
                (e) => e.course.id === course);
        }
    }

    private selectAssignment(labIdString: string) {
        const labId = parseInt(labIdString, 10);
        if (this.selectedCourse && !isNaN(labId)) {
            // TODO: Be carefull not to return anything that sould not be able to be returned
            this.selectedAssignment = this.selectedCourse.assignments.find(
                (e) => e.assignment.id === labId,
            );
        }
    }

    private handleClick(link: ILink) {
        if (link.uri) {
            this.navMan.navigateTo(link.uri);
        }
    }
}
