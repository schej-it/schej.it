import 'package:auto_route/auto_route.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/models/event.dart';
import 'package:flutter_app/pages/events/create_event.dart';
import 'package:flutter_app/router/app_router.gr.dart';
import 'package:flutter_app/utils.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';
import '../../components/event_card.dart';
import '../../constants/colors.dart';
import '../../constants/fonts.dart';

class EventsPage extends StatefulWidget {
  const EventsPage({Key? key}) : super(key: key);

  @override
  State<EventsPage> createState() => _EventsPageState();
}

class _EventsPageState extends State<EventsPage> {
  // Sample events data.
  final List<Event> events = List<Event>.generate(
      5,
      (i) => Event(
          id: '39ajfa',
          ownerId: 'tonyxin',
          name: 'Fun times!',
          startDate: getDateWithTime(
              DateTime.now().subtract(const Duration(days: 1)), 0),
          endDate: getDateWithTime(DateTime.now(), 0),
          responses: List<String>.generate(5, (i) => "Tommy")));

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SchejAppBar(titleString: 'My events', isRoot: true),
      floatingActionButton: FloatingActionButton(
        heroTag: 'eventsFab',
        onPressed: () => Navigator.of(context).push(_createRoute()),
        backgroundColor: SchejColors.darkGreen,
        child: const Icon(MdiIcons.plus),
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(25.0),
          child: Center(
            child: Column(
              children: <Widget>[
                _buildEventsSection('Events I created'),
                const SizedBox(height: 20),
                _buildEventsSection('Events I joined'),
              ],
            ),
          ),
        ),
      ),
    );
  }

  // Builds the events section headers with a see all option.
  Widget _buildEventsSection(String title) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.only(bottom: 13.0),
          child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: <Widget>[
                Text(title,
                    style: SchejFonts.subtitle
                        .copyWith(color: SchejColors.darkGreen)),
                GestureDetector(
                  onTap: () {
                    AutoRouter.of(context)
                        .push(EventPageRoute(eventId: '1234'));
                  },
                  child: const Text('See all', style: SchejFonts.body),
                )
              ]),
        ),
        ListView.separated(
            shrinkWrap: true,
            itemCount: events.length,
            itemBuilder: (context, index) {
              return EventCard(event: events[index]);
            },
            separatorBuilder: (context, index) => const SizedBox(
                  height: 10,
                )),
      ],
    );
  }
}

Route _createRoute() {
  return PageRouteBuilder(
    pageBuilder: (context, animation, secondaryAnimation) =>
        const CreateEvent(),
    transitionsBuilder: (context, animation, secondaryAnimation, child) {
      const begin = Offset(0.0, 1.0);
      const end = Offset.zero;
      const curve = Curves.ease;

      var tween = Tween(begin: begin, end: end).chain(CurveTween(curve: curve));

      return SlideTransition(
        position: animation.drive(tween),
        child: child,
      );
    },
  );
}
