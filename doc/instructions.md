# Usage

## General Use

  * GET `/` Print Instructions (this page)
  * GET `/target/:groupID` Print the current target for the given group ID; an optional `date` parameter may be passed to resolve the schedule for that date instead of now.

## Groups

A `group` has the data structure:
```json
			{
				"id": "ID of group",
				"name": "name/label of group",
				"timezone": "Time zone, of the form US/Eastern or America/New York",
			}
```

  * **GET** `/group/:groupID` Print the group identified by groupID
  * **POST** `/group` Add a group.

## Import

There are two types of CSV import:  "days" and "dates".  "days" imports a default schedule, based
on the provided generic days of the week.  "dates" imports schedules for specific dates.  If there
exists a "dates" schedule for any given time, it is used in preference to the "days" schedule.

NOTE: if the target is blank (""), the row will be ignored.

A "days" schedule is a CSV file with no field headers and columns of the form:
```
   "Group ID","Day of the Week","Start Time (HH:MM)","Stop Time (HH:MM)","Target phone number"
```
_(Day of the Week can be one- or three-letter abbreviations or the full weekday name: 'M', 'Mon', 'Monday')_

A "dates" schedule is a CSV file with no field headers and columns of the form:
```
   "Group ID","Date (YYYY-MM-DD)","Start Time (HH:MM)","Stop Time (HH:MM)","Target phone number"
```

  * **POST** `/sched/import/days` Add a days (generic weekly) schedule.
  * **POST** `/sched/import/dates` Add a dates (specific dates) schedule.
