import 'package:auto_route/auto_route.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/event.dart';
import 'package:flutter_app/router/app_router.gr.dart';
import 'package:intl/intl.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

// Displays a card that shows an event.
class EventCard extends StatelessWidget {
  final Event event;

  const EventCard({Key? key, required this.event}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
        onTap: () {
          AutoRouter.of(context).push(EventPageRoute(eventId: event.id));
        },
        child: DefaultTextStyle(
          style: SchejFonts.body,
          child: Container(
            padding: const EdgeInsets.all(10.0),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: const BorderRadius.all(Radius.circular(10.0)),
              boxShadow: [
                BoxShadow(
                  color: Colors.grey.withOpacity(0.3),
                  spreadRadius: 2,
                  blurRadius: 2,
                  offset: const Offset(0, 1), // changes position of shadow
                ),
              ],
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    _buildTitleDate(),
                    Padding(
                      padding: const EdgeInsets.only(left: 8.0),
                      child: _buildResponses(),
                    ),
                  ],
                ),
                Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    IconButton(
                      icon: const Icon(MdiIcons.chevronRight,
                          color: SchejColors.darkGray, size: 35.0),
                      splashRadius: 15,
                      onPressed: () => {},
                    ),
                  ],
                ),
              ],
            ),
          ),
        ));
  }

  Column _buildTitleDate() {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(event.name),
        Padding(
          padding: const EdgeInsets.only(top: 5.0),
          child: Text(DateFormat('Md').format(event.startDate),
              style: SchejFonts.small.copyWith(color: SchejColors.darkGray)),
        ),
      ],
    );
  }

  // Builds responses icon section. Separated for readability.
  Container _buildResponses() {
    return Container(
      padding: const EdgeInsets.all(5.0),
      decoration: const BoxDecoration(
          color: SchejColors.offWhite,
          borderRadius: BorderRadius.all(Radius.circular(10.0))),
      child: Row(
        children: [
          const Icon(MdiIcons.accountMultiple, size: 13.0),
          const SizedBox(width: 2),
          Text("${event.responses.length}", style: SchejFonts.small)
        ],
      ),
    );
  }
}
