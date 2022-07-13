import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/fonts.dart';

class Calendar extends StatefulWidget {
  const Calendar({Key? key}) : super(key: key);

  @override
  State<Calendar> createState() => _CalendarState();
}

class _CalendarState extends State<Calendar> {
  final double timeColWidth = 50;
  final double timeRowHeight = 45;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        _buildDaySection(),
        const Divider(color: SchejColors.darkGray),
        Expanded(child: _buildCalendarSection()),
      ],
    );
  }

  Widget _buildDaySection() {
    return Row(
      // mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        SizedBox(width: timeColWidth),
        Expanded(child: _buildDay('Tue', 5)),
        Expanded(child: _buildDay('Wed', 6)),
        Expanded(child: _buildDay('Thu', 7)),
      ],
    );
  }

  Widget _buildDay(String dayText, int dateNum) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Text(dayText, style: SchejFonts.body),
        Text(dateNum.toString(), style: SchejFonts.header),
      ],
    );
  }

  Widget _buildCalendarSection() {
    final timeStrings = <String>[];
    for (int i = 1; i < 24; ++i) {
      String timeText;
      if (i < 12) {
        timeText = '$i AM';
      } else if (i == 12) {
        timeText = '12 PM';
      } else {
        timeText = '${i - 12} PM';
      }
      timeStrings.add(timeText);
    }

    // return const Text('test');
    return ListView.builder(
      scrollDirection: Axis.vertical,
      itemCount: timeStrings.length,
      controller: ScrollController(initialScrollOffset: 8.25*timeRowHeight),
      itemBuilder: (BuildContext context, int index) {
        return SizedBox(
          height: timeRowHeight,
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              SizedBox(
                width: timeColWidth,
                child: Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.only(right: 8.0),
                      child: Align(
                        alignment: Alignment.centerRight,
                        child: Text(
                          timeStrings[index],
                          style: SchejFonts.body.copyWith(color: SchejColors.darkGray),
                          textAlign: TextAlign.right,
                        ),
                      ),
                    ),
                    const Spacer(),
                    const SizedBox(width: 5, child: Divider(color: SchejColors.lightGray)),
                  ],
                ),
              ),
              const VerticalDivider(width: 1, color: SchejColors.lightGray),
              const Expanded(child: Divider(color: SchejColors.lightGray)),
              const VerticalDivider(width: 1, color: SchejColors.lightGray),
              const Expanded(child: Divider(color: SchejColors.lightGray)),
              const VerticalDivider(width: 1, color: SchejColors.lightGray),
              const Expanded(child: Divider(color: SchejColors.lightGray)),
            ],
          ),
        );
      },
    );
  }
}
