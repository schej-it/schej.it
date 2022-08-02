import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';

// Widget containing chips of all the people added to the current schej comparison
// as well as a text field at the right to search for more people
class CompareSchejTextField extends StatefulWidget {
  final FocusNode focusNode;
  final TextEditingController controller;
  final ValueChanged<String>? onChanged;
  final Set<String>? addedUsers;
  final bool includeSelf;

  const CompareSchejTextField({
    Key? key,
    required this.focusNode,
    required this.controller,
    this.onChanged,
    this.addedUsers,
    this.includeSelf = true,
  }) : super(key: key);

  @override
  State<CompareSchejTextField> createState() => _CompareSchejTextFieldState();
}

class _CompareSchejTextFieldState extends State<CompareSchejTextField> {
  // Variables
  bool _focused = false;

  @override
  void initState() {
    super.initState();
    widget.focusNode.addListener(setFocused);
  }

  @override
  void dispose() {
    widget.focusNode.removeListener(setFocused);
    super.dispose();
  }

  void setFocused() {
    setState(() {
      _focused = widget.focusNode.hasFocus;
    });
  }

  @override
  Widget build(BuildContext context) {
    final tag = Container(
      decoration: BoxDecoration(
        borderRadius: SchejConstants.borderRadius,
        color: SchejColors.green,
      ),
      margin: const EdgeInsets.symmetric(horizontal: 5),
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      child: Row(
        children: [
          InkWell(
            child: Text('Jonathan Liu',
                style: SchejFonts.subtitle.copyWith(color: SchejColors.white)),
            onTap: () {
              // print('$tag selected!');
            },
          ),
          const SizedBox(width: 4),
          InkWell(
            child: const Icon(
              Icons.cancel,
              size: 14,
              color: SchejColors.white,
            ),
            onTap: () {
              // onDeleteTag(tag);
            },
          ),
        ],
      ),
    );

    final textField = ConstrainedBox(
      constraints: const BoxConstraints(maxWidth: 255),
      child: TextField(
        focusNode: widget.focusNode,
        controller: widget.controller,
        onChanged: widget.onChanged,
        autocorrect: false,
        decoration: InputDecoration(
          hintText: 'Type a name...',
          contentPadding: const EdgeInsets.symmetric(horizontal: 12),
          enabledBorder: OutlineInputBorder(
            borderRadius: SchejConstants.borderRadius,
            borderSide: BorderSide.none,
          ),
          focusedBorder: OutlineInputBorder(
            borderRadius: SchejConstants.borderRadius,
            borderSide: BorderSide.none,
          ),
        ),
        style: SchejFonts.subtitle,
      ),
    );

    return Container(
      width: double.infinity,
      decoration: BoxDecoration(
        color: _focused ? SchejColors.white : SchejColors.offWhite,
        borderRadius: SchejConstants.borderRadius,
        border: Border.all(width: 1, color: SchejColors.lightGray),
      ),
      margin: const EdgeInsets.only(bottom: 20),
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          mainAxisSize: MainAxisSize.max,
          children: [
            // tag,
            tag,
            textField,
          ],
        ),
      ),
    );
  }
}
