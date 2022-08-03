import 'package:animations/animations.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/components/calendar.dart';
import 'package:flutter_app/components/calendar_view_selector.dart';
import 'package:flutter_app/components/friends/compare_schej_text_field.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/models/calendar_event.dart';
import 'package:flutter_app/pages/friends/compare_schej_dialog.dart';
import 'package:flutter_app/utils.dart';

class FriendSchejPage extends StatefulWidget {
  final String friendId;

  const FriendSchejPage({Key? key, required this.friendId}) : super(key: key);

  @override
  State<FriendSchejPage> createState() => _FriendSchejPageState();
}

class _FriendSchejPageState extends State<FriendSchejPage> {
  // Calendar variables
  int _daysVisible = 3;
  DateTime _selectedDay = getDateWithTime(DateTime.now(), 0);

  // Compare schej dialog variables
  final Set<String> _addedUsers = <String>{};
  bool _includeSelf = true;

  @override
  void initState() {
    super.initState();

    _addedUsers.add(widget.friendId);
  }

  @override
  Widget build(BuildContext context) {
    final testCalendarEvents = CalendarEvents(
      events: [
        CalendarEvent(
          title: 'Event',
          startDate: getDateWithTime(DateTime.now(), 9.5),
          endDate: getDateWithTime(DateTime.now(), 12),
        ),
        CalendarEvent(
          title: 'Introduction to Failure Analysis',
          startDate: getDateWithTime(DateTime.now(), 13),
          endDate: getDateWithTime(DateTime.now(), 14.5),
        ),
        CalendarEvent(
          title: 'War',
          startDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 1)), 15),
          endDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 1)), 20),
        ),
      ],
    );

    return Scaffold(
      appBar: SchejAppBar(
        centerTitle: true,
        underline: false,
        titleString: 'Compare schej',
        actions: [
          CalendarViewSelector(
            onSelected: (int value) => setState(() {
              _daysVisible = value;
            }),
          ),
        ],
      ),
      body: Container(
        color: SchejColors.white,
        child: Stack(
          children: [
            Positioned.fill(
              child: Calendar(
                calendarEvents: testCalendarEvents,
                daysVisible: _daysVisible,
                selectedDay: _selectedDay,
                onDaySelected: (selectedDay) => setState(() {
                  _selectedDay = selectedDay;
                }),
              ),
            ),
            _buildTextField(),
          ],
        ),
      ),
    );
  }

  Widget _buildTextField() {
    return Positioned(
      bottom: 0,
      width: MediaQuery.of(context).size.width,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 25, horizontal: 18),
        child: OpenContainer(
          closedShape: RoundedRectangleBorder(
            borderRadius: SchejConstants.borderRadius,
          ),
          closedBuilder: (context, openContainer) {
            return StatefulBuilder(
              builder: (context, setState) => FocusScope(
                child: Focus(
                  onFocusChange: (focus) {
                    if (focus) {
                      openContainer();
                    }
                  },
                  child: CompareSchejTextField(
                    addedUsers: _addedUsers,
                    onRemoveUser: (userId) => setState(() {
                      _addedUsers.remove(userId);
                    }),
                    includeSelf: _includeSelf,
                  ),
                ),
              ),
            );
          },
          openBuilder: (context, closeContainer) {
            return StatefulBuilder(
              builder: (context, setState) => CompareSchejDialog(
                addedUsers: _addedUsers,
                onAddUser: ((userId, added) {
                  setState(() {
                    if (added) {
                      _addedUsers.add(userId);
                    } else {
                      _addedUsers.remove(userId);
                    }
                  });
                }),
                includeSelf: _includeSelf,
                onUpdateIncludeSelf: (value) =>
                    setState(() => _includeSelf = value),
                onClose: closeContainer,
              ),
            );
          },
        ),
      ),
    );
  }
}
