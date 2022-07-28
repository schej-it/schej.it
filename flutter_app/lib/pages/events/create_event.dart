import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

class CreateEvent extends StatefulWidget {
  const CreateEvent({Key? key}) : super(key: key);

  @override
  State<CreateEvent> createState() => _CreateEventState();
}

class _CreateEventState extends State<CreateEvent> {
  String startTime = 'One';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          children: [
            _buildHeaderSection(context),
            Container(
                padding: const EdgeInsets.all(10.0),
                child: Column(
                  children: [
                    _buildEventName(),
                    const SizedBox(height: 20.0),
                    _buildTimesToMeet(),
                  ],
                )),
          ],
        ),
      ),
    );
  }
}

Widget _buildHeaderSection(BuildContext context) {
  return Stack(
    children: [
      Positioned(
        top: 0,
        right: 10,
        child: IconButton(
          icon:
              const Icon(MdiIcons.close, color: SchejColors.black, size: 30.0),
          splashRadius: 15,
          onPressed: () => Navigator.pop(context),
        ),
      ),
      Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: const [
          Padding(
            padding: EdgeInsets.only(top: 8.0),
            child: Text('New event', style: SchejFonts.header),
          ),
        ],
      ),
    ],
  );
}

Widget _buildEventName() {
  return TextFormField(
    decoration: const InputDecoration(
      hintText: 'Name Your Event',
      hintStyle: SchejFonts.header,
      enabledBorder: UnderlineInputBorder(
        borderSide: BorderSide(color: SchejColors.lightGray),
      ),
      focusedBorder: UnderlineInputBorder(
        borderSide: BorderSide(color: SchejColors.lightGray),
      ),
    ),
    style: SchejFonts.header,
    cursorColor: SchejColors.darkGreen,
  );
}

Widget _buildTimesToMeet() {
  return Column(
    crossAxisAlignment: CrossAxisAlignment.start,
    children: [
      const Text('What times would you like to meet between?',
          style: SchejFonts.subtitle),
      const SizedBox(height: 10.0),
      Container(
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
        child: Column(
          children: [
            _buildTimeInput('Start'),
            _buildTimeInput('End'),
          ],
        ),
      ),
    ],
  );
}

Widget _buildTimeInput(String text) {
  return Row(
    mainAxisAlignment: MainAxisAlignment.spaceBetween,
    children: [
      const Text('hi'),
      DropdownButton<String>(
        value: 'One',
        icon: const Icon(Icons.arrow_downward),
        elevation: 16,
        style: const TextStyle(color: Colors.deepPurple),
        underline: Container(
          height: 2,
          color: Colors.deepPurpleAccent,
        ),
        onChanged: (String? newValue) {
          // print(newValue);
        },
        items: <String>['One', 'Two', 'Free', 'Four']
            .map<DropdownMenuItem<String>>((String value) {
          return DropdownMenuItem<String>(
            value: value,
            child: Text(value),
          );
        }).toList(),
      ),
    ],
  );
}
