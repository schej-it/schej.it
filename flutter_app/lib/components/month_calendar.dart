import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:table_calendar/table_calendar.dart';

import '../constants/colors.dart';
import '../constants/fonts.dart';

// MonthCalendar displays a month view of the calendar, and allows for date
// selection
class MonthCalendar extends StatefulWidget {
  final DateTime selectedDay;
  final void Function(DateTime) onDaySelected;

  const MonthCalendar({
    Key? key,
    required this.selectedDay,
    required this.onDaySelected,
  }) : super(key: key);

  @override
  State<MonthCalendar> createState() => _MonthCalendarState();
}

class _MonthCalendarState extends State<MonthCalendar> {
  @override
  Widget build(BuildContext context) {
    return TableCalendar(
      focusedDay: widget.selectedDay,
      firstDay: DateTime.now().subtract(const Duration(days: 365)),
      lastDay: DateTime.now().add(const Duration(days: 365)),
      headerVisible: false,
      daysOfWeekHeight: 30,
      rowHeight: 30,
      daysOfWeekStyle: DaysOfWeekStyle(
        dowTextFormatter: (date, locale) => DateFormat('E').format(date)[0],
        weekdayStyle: SchejFonts.smallBold,
        weekendStyle: SchejFonts.smallBold,
      ),
      calendarStyle: CalendarStyle(
        defaultTextStyle: SchejFonts.small,
        disabledTextStyle: SchejFonts.small,
        todayTextStyle: SchejFonts.small,
        todayDecoration: BoxDecoration(
          color: SchejColors.darkGreen.withOpacity(0.25),
          // border: Border.all(color: SchejColors.gray, width: 1),
          shape: BoxShape.circle,
        ),
        weekendTextStyle: SchejFonts.small,
        // outsideTextStyle: SchejFonts.small.copyWith(color: SchejColors.gray),
        outsideDaysVisible: false,
        selectedTextStyle: SchejFonts.small.copyWith(color: SchejColors.white),
        selectedDecoration: const BoxDecoration(
          color: SchejColors.darkGreen,
          shape: BoxShape.circle,
        ),
      ),
      selectedDayPredicate: (day) {
        return isSameDay(widget.selectedDay, day);
      },
      onDaySelected: (selectedDay, focusedDay) {
        if (selectedDay != widget.selectedDay) {
          widget.onDaySelected(selectedDay);
        }
      },
      onPageChanged: (focusedDay) {
        if (widget.selectedDay != focusedDay) {
          widget.onDaySelected(focusedDay);
        }
      },
    );
  }
}
