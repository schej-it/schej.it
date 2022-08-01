import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/components/friends/compare_schej_card.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';

class CompareSchejDialog extends StatefulWidget {
  final VoidCallback onClose;

  const CompareSchejDialog({Key? key, required this.onClose}) : super(key: key);

  @override
  State<CompareSchejDialog> createState() => _CompareSchejDialogState();
}

class _CompareSchejDialogState extends State<CompareSchejDialog> {
  final _focusNode = FocusNode();

  var results = [
    {'name': 'Winston Tilton', 'added': false},
    {'name': 'Samantha Hutchinson', 'added': false},
    {'name': 'Tyler Smithson', 'added': true},
    {'name': 'Arthi Singh', 'added': false},
  ];

  @override
  initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _focusNode.requestFocus();
    });
  }

  @override
  void dispose() {
    _focusNode.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SchejAppBar(
        isRoot: true,
        underline: false,
        centerTitle: true,
        titleString: 'Compare schej',
        actions: [
          IconButton(onPressed: widget.onClose, icon: const Icon(Icons.close)),
        ],
      ),
      body: Container(
        color: SchejColors.white,
        padding: SchejConstants.pagePadding,
        child: Column(
          children: [
            _buildSearchTextField(),
            Expanded(child: _buildResults()),
          ],
        ),
      ),
    );
  }

  Widget _buildSearchTextField() {
    final textField = Padding(
      padding: const EdgeInsets.only(bottom: 20),
      child: TextField(
        focusNode: _focusNode,
        autocorrect: false,
        decoration: const InputDecoration(
          hintText: 'Type a name...',
          prefixIcon: Icon(Icons.search),
        ),
        style: SchejFonts.subtitle,
      ),
    );

    return textField;
  }

  Widget _buildResults() {
    return ListView.builder(
      itemCount: results.length,
      itemBuilder: (context, index) {
        final result = results[index];
        return Padding(
          padding: const EdgeInsets.only(bottom: 10),
          child: CompareSchejCard(
            name: result['name'] as String,
            added: result['added'] as bool,
            onToggle: (bool value) {
              setState(() {
                results[index]['added'] = value;
              });
            },
          ),
        );
      },
    );
  }
}
